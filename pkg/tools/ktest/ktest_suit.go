// Copyright 2021 The KCL Authors. All rights reserved.

package ktest

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/go-cmp/cmp"
	"gopkg.in/yaml.v3"

	"kusionstack.io/kclvm-go/pkg/compiler/parser"
	"kusionstack.io/kclvm-go/pkg/kcl"
	"kusionstack.io/kclvm-go/pkg/service"
	"kusionstack.io/kclvm-go/pkg/spec/gpyrpc"
)

const __kcl_test_main = "__kcl_test_main.k"

type kTestSuit struct {
	opt Options

	workDir string
	pkgpath string

	kFiles        []string
	kSettingFiles []string
	kPluginFiles  []string
	kTestFiles    []string

	hasGoldenFile bool
	goldenStdout  string // stdout.golden
	goldenStderr  string // stderr.golden

	tSchemaNames   []string
	tSchemaInfoMap map[string]*kTestSchemaInfo

	skipMsg string
	skip    bool

	err error
}

type kTestSchemaInfo struct {
	Args []*gpyrpc.CmdArgSpec

	SettingFile      string
	SettingFile_Args []*gpyrpc.CmdArgSpec
}

func loadKTestSuit(workDir string, opt Options) (*kTestSuit, error) {
	kFiles, kSettingFiles, kPluginFiles, kTestFiles := readKclFiles(workDir)

	klog.Debugf("kFiles = %v\n", kFiles)
	klog.Debugf("kSettingFiles = %v\n", kSettingFiles)
	klog.Debugf("kPluginFiles = %v\n", kPluginFiles)
	klog.Debugf("kTestFiles = %v\n", kTestFiles)

	t := &kTestSuit{
		opt: opt,

		workDir: workDir,
		pkgpath: workDir,

		kFiles:        kFiles,
		kSettingFiles: kSettingFiles,
		kPluginFiles:  kPluginFiles,
		kTestFiles:    kTestFiles,

		tSchemaInfoMap: make(map[string]*kTestSchemaInfo),
	}

	// stdout.golden or stdout.golden.py
	if _, err := os.Stat(filepath.Join(workDir, "stdout.golden")); err == nil {
		d, _ := os.ReadFile(filepath.Join(workDir, "stdout.golden"))
		t.goldenStdout = string(d)
		t.hasGoldenFile = true
	}

	// stderr.golden or stderr.golden.py
	if _, err := os.Stat(filepath.Join(workDir, "stderr.golden")); err == nil {
		d, _ := os.ReadFile(filepath.Join(workDir, "stderr.golden"))
		t.goldenStderr = strings.ReplaceAll(string(d), "${PWD}", workDir)
		t.hasGoldenFile = true
	}

	for _, kTestFile := range kTestFiles {
		f, err := parser.ParseFile(kTestFile, nil)
		if err != nil {
			t.err = err
			return t, err
		}

		for _, tSchemaName := range getTestSchemaNameList(f) {
			if opt.shouldRun(tSchemaName) {
				tSchemaInfo, err := getTestSchemaInfo(workDir, f, tSchemaName)
				if err != nil {
					t.err = err
					return t, err
				}

				t.tSchemaNames = append(t.tSchemaNames, tSchemaName)
				t.tSchemaInfoMap[tSchemaName] = tSchemaInfo

				klog.Debugf("tSchemaName = %s\n", tSchemaName)
				klog.Debugf("tSchemaInfo = %v\n", tSchemaInfo)
			}
		}
	}

	t.skipMsg, t.skip = t.shouldSkip()

	return t, nil
}

func (p *kTestSuit) RunTest() error {
	if p.skip {
		if !p.opt.QuietMode {
			if p.skipMsg != "" {
				fmt.Printf("skip %v [%s]\n", p.pkgpath, p.skipMsg)
			} else {
				fmt.Printf("skip %v\n", p.pkgpath)
			}
		}
		return nil
	}
	if p.err != nil {
		if !p.opt.QuietMode {
			fmt.Printf("FAIL %v err: %v", p.pkgpath, p.err)
		}
		return p.err
	}
	if len(p.kFiles) == 0 {
		if len(p.kTestFiles) == 0 {
			if !p.opt.QuietMode {
				fmt.Printf("?    %v [no test files]", p.pkgpath)
			}
			return nil
		}
		if len(p.tSchemaNames) == 0 {
			if !p.opt.QuietMode {
				fmt.Printf("ok   %v [no tests to run]\n", p.pkgpath)
			}
			return nil
		}
	}

	// change dir
	if wd, err := os.Getwd(); err == nil {
		err := os.Chdir(p.workDir)
		defer os.Chdir(wd)

		if err != nil {
			return errors.New("chdir failed: " + p.workDir)
		}
	}

	startTime := time.Now()

	// kcl-test xxx_test.k
	ok, msg := p.runTest()
	if !ok {
		if !p.opt.QuietMode {
			fmt.Printf("FAIL %v [%v]\n", p.pkgpath, time.Now().Sub(startTime))
		}
		if msg = strings.TrimSpace(msg); msg != "" {
			if !p.opt.QuietMode {
				fmt.Printf("%s\n", msg)
			}
		}
		return errors.New(msg)
	}

	// kcl main.k
	ok2, msg2 := p.runMainK()
	if !ok2 {
		if !p.opt.QuietMode {
			fmt.Printf("FAIL %v [%v]\n", p.pkgpath, time.Now().Sub(startTime))
		}
		if msg2 = strings.TrimSpace(msg2); msg2 != "" {
			if !p.opt.QuietMode {
				fmt.Printf("%s\n", msg2)
			}
		}
		return errors.New(msg2)
	}

	// OK
	if !p.opt.QuietMode {
		fmt.Printf("ok   %v [%v]\n", p.pkgpath, time.Now().Sub(startTime))
		if msg = strings.TrimSpace(msg); msg != "" {
			fmt.Printf("%s\n", msg)
		}
		if msg2 = strings.TrimSpace(msg2); msg2 != "" {
			fmt.Printf("%s\n", msg2)
		}
	}
	return nil
}

func (p *kTestSuit) runTest() (ok bool, msg string) {
	if len(p.tSchemaNames) == 0 {
		return true, "" // skip
	}

	var buf bytes.Buffer

	if !p.opt.Debug {
		defer os.Remove(__kcl_test_main)
	}

	p.genTestMainFile()

	var all_k_files []string
	all_k_files = append(all_k_files, p.kFiles...)
	all_k_files = append(all_k_files, p.kTestFiles...)
	all_k_files = append(all_k_files, filepath.Join(p.workDir, __kcl_test_main))

	client := service.NewKclvmServiceClient()

	// test plugin
	if len(p.kPluginFiles) > 0 {
		pluginRootBackup := os.Getenv("KCL_PLUGINS_ROOT")
		defer func() {
			os.Setenv("KCL_PLUGINS_ROOT", pluginRootBackup)
			client.ResetPlugin(&gpyrpc.ResetPlugin_Args{PluginRoot: pluginRootBackup})
		}()

		if _, err := os.Stat(filepath.Join(p.workDir, "plugin.py")); err == nil {
			os.Setenv("KCL_PLUGINS_ROOT", filepath.Dir(p.workDir))
		}

		klog.Debugf("KCL_PLUGINS_ROOT = %s\n", os.Getenv("KCL_PLUGINS_ROOT"))

		_, err := client.ResetPlugin(&gpyrpc.ResetPlugin_Args{
			PluginRoot: os.Getenv("KCL_PLUGINS_ROOT"),
		})
		if err != nil {
			fmt.Fprintf(&buf, "---- reset_plugin failed\n")
			fmt.Fprintf(&buf, "%s\n", withLinePrefix(err.Error(), "     "))
			return false, buf.String()
		}
	}

	// only try compile
	_, err := client.ExecProgram(&gpyrpc.ExecProgram_Args{
		WorkDir:       p.workDir,
		KFilenameList: all_k_files,
		Args: []*gpyrpc.CmdArgSpec{
			{Name: "__kcl_test_run", Value: "___test_schema_@@@__"},
			{Name: "__kcl_test_debug", Value: fmt.Sprint(p.opt.Debug)},
		},
		Overrides:         []*gpyrpc.CmdOverrideSpec{},
		DisableYamlResult: true,
	})
	if err != nil {
		fmt.Fprintf(&buf, "---- compile failed\n")
		fmt.Fprintf(&buf, "%s\n", withLinePrefix(err.Error(), "     "))
		return false, buf.String()
	}

	// run test list
	var allOK = true
	for _, testSchemaName := range p.tSchemaNames {
		args := &gpyrpc.ExecProgram_Args{
			WorkDir:       p.workDir,
			KFilenameList: all_k_files,
			Args: []*gpyrpc.CmdArgSpec{
				{Name: "__kcl_test_run", Value: testSchemaName},
				{Name: "__kcl_test_debug", Value: fmt.Sprint(p.opt.Debug)},
			},
			Overrides:         []*gpyrpc.CmdOverrideSpec{},
			DisableYamlResult: true,
		}

		tSchemaInfo := p.tSchemaInfoMap[testSchemaName]
		args.Args = append(args.Args, tSchemaInfo.Args...)

		startTime := time.Now()
		_, err := client.ExecProgram(args)
		timeUsed := time.Now().Sub(startTime)
		if err != nil {
			allOK = false
			errMsg := withLinePrefix(err.Error(), "     ")
			fmt.Fprintf(&buf, "---- <%s> failed [%v]\n", testSchemaName, timeUsed)
			fmt.Fprintf(&buf, "%s\n", errMsg)
			continue
		}
		if p.opt.Verbose {
			fmt.Fprintf(&buf, "---- <%s> success [%v]\n", testSchemaName, timeUsed)
			continue
		}
	}

	return allOK, buf.String()
}

func (p *kTestSuit) runMainK() (ok bool, msg string) {
	if !(p.hasMainK() && p.hasGoldenFile) || len(p.kPluginFiles) != 0 {
		return true, ""
	}

	var paths []string
	paths = append(paths, p.kFiles...)
	paths = append(paths, p.kSettingFiles...)

	results, err := kcl.RunFiles(paths, kcl.WithWorkDir(p.workDir))
	if err != nil {
		return p.checkGoldenStderr(err)
	}

	return p.checkGoldenStdout(results)
}

func (p *kTestSuit) hasMainK() bool {
	// TODO: fix settings.yaml
	if fi, _ := os.Stat(filepath.Join(p.workDir, "settings.yaml")); fi != nil {
		return false
	}
	for _, s := range p.kFiles {
		if filepath.Base(s) == "main.k" {
			p.kFiles = []string{s}
			return true
		}
	}
	return false
}

func (p *kTestSuit) shouldSkip() (skipMsg string, ok bool) {
	if len(p.kTestFiles) == 0 {
		return "no test", true
	}
	for _, s := range p.kFiles {
		data, _ := os.ReadFile(s)
		if strings.Contains(string(data), "# kcl-test: ignore") {
			return "# kcl-test: ignore", true
		}
	}
	for _, s := range p.kPluginFiles {
		data, _ := os.ReadFile(s)
		if strings.Contains(string(data), "# kcl-test: ignore") {
			return "# kcl-test: ignore", true
		}
	}
	return "", false
}

func (p *kTestSuit) genTestMainFile() {
	test_main_k_code := genTestMainFile(p.tSchemaNames)
	os.WriteFile(__kcl_test_main, []byte(test_main_k_code), 0666)
}

func (p *kTestSuit) checkGoldenStderr(err error) (ok bool, msg string) {
	want, got := p.goldenStderr, err.Error()
	if diff := cmp.Diff(want, got); diff != "" {
		return false, fmt.Sprintf("golden stderr mismatch (-want +got):\n%s", diff)
	} else {
		return true, ""
	}
}
func (p *kTestSuit) checkGoldenStdout(results *kcl.KCLResultList) (ok bool, msg string) {
	want, err := goldenYamlString(p.goldenStdout)
	if err != nil {
		return false, fmt.Sprintf("want: unsupport format: %v", err)
	}
	got, err := goldenYamlString(results.First().YAMLString())
	if err != nil {
		return false, fmt.Sprintf("got: unsupport format: %v", err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		return false, fmt.Sprintf("golden stdout mismatch (-want +got):\n%s\nwant:%s\ngot:%s", diff, want, got)
	} else {
		return true, ""
	}
}

func goldenYamlString(s string) (string, error) {
	// decode yaml (keep equal to json)
	var m map[string]interface{}
	err := yaml.Unmarshal([]byte(s), &m)
	if err != nil {
		return "", err
	}

	// yaml -> json -> yaml
	{
		jsonString, err := json.Marshal(m)
		if err != nil {
			return "", err
		}
		if err = yaml.Unmarshal(jsonString, &m); err != nil {
			return "", err
		}
	}

	d, err := yaml.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(d), nil
}
