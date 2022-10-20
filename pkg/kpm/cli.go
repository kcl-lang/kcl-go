package kpm

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/user"
	"strings"
)

func CLI(args ...string) {
	if len(args) < 1 {
		println(CliHelp)
		return
	}
	err := CliSetup()

	if err != nil {
		return
	}
	switch args[0] {
	case "init":
		if len(args) != 2 {
			println(CliInitHelp)
			return
		}
		err = CliInit(args[1])
		if err != nil {
			println(err.Error())
			return
		}

	case "add":
		if len(args) < 2 {
			println(CliAddHelp)

			return
		}
		err = CliAdd(args[1:]...)
		if err != nil {
			println(err.Error())
			return
		}

	case "del":
		if len(args) < 2 {
			println(CliDelHelp)
			return
		}
		err = CliDel(args[1:]...)
		if err != nil {
			println(err.Error())
			return
		}

	default:
		println(CliNotFound)
		println(CliHelp)
		//弹出使用方法
	}
}

// CliSetup 加载环境变量，初始化目录与设置
func CliSetup() error {
	var err error
	pwd, err = os.Getwd()
	if err != nil {
		return nil
	}
	//加载环境变量
	if tmp := os.Getenv("KPM_ROOT"); tmp == "" {
		home := ""
		u, err := user.Current()
		if err != nil {
			if tmphome := os.Getenv("HOME"); tmphome != "" {
				home = tmphome
			} else {
				return nil
			}
		}
		home = u.HomeDir
		KPM_ROOT = home + Separator + "kpm"
	}
	if tmp := os.Getenv("KPM_SERVER_ADDR"); tmp != "" {
		KPM_SERVER_ADDR = tmp
	}
	parse, err := url.Parse(KPM_SERVER_ADDR)
	if err != nil {
		return err
	}
	KPM_SERVER_ADDR_PATH = parse.Host

	//初始化目录信息
	err = KeepDirExists(KPM_ROOT,
		KPM_ROOT+Separator+"registry",
		KPM_ROOT+Separator+"registry"+Separator+KPM_SERVER_ADDR_PATH,
		KPM_ROOT+Separator+"registry"+Separator+KPM_SERVER_ADDR_PATH+Separator+"kcl_modules",
		KPM_ROOT+Separator+"registry"+Separator+KPM_SERVER_ADDR_PATH+Separator+"tag",
		KPM_ROOT+Separator+"registry"+Separator+KPM_SERVER_ADDR_PATH+Separator+"metadata",
		KPM_ROOT+Separator+"git",
		KPM_ROOT+Separator+"git"+Separator+"kcl_modules",
		KPM_ROOT+Separator+"git"+Separator+"metadata",
		KPM_ROOT+Separator+"store",
		KPM_ROOT+Separator+"store"+Separator+"v1",
		KPM_ROOT+Separator+"store"+Separator+"v1"+Separator+"files",
	)
	if err != nil {
		println("setup fail,", err.Error())
		return err
	}
	for i := 0; i < len(hextable); i++ {
		for j := 0; j < len(hextable); j++ {
			err = KeepDirExists(KPM_ROOT + Separator + "store" + Separator + "v1" + Separator + "files" +
				Separator + string(hextable[i]) + string(hextable[j]))
			if err != nil {
				return err
			}
		}
	}
	version, err := GetKclvmMinVersion()
	if err == nil {
		KclvmMinVersion = "v" + version
	}

	return nil
}

// CliAdd 添加包，检查vm版本，如果比当前版本大，则失败，只负责链接或者复制
func CliAdd(args ...string) error {
	//flag_global := false
	flag_git := false
	//flag_internal := false
	var pkgvs []string
	//var pkgs []Require
	for i := 0; i < len(args); i++ {
		if strings.HasPrefix(args[i], "-") {
			switch args[i] {
			//case "-g":
			//	flag_global = true
			case "-git":
				flag_git = true
				//case "--internal":
				//	flag_internal = true
			}
		} else {
			pkgvs = args[i:]
			break
		}
	}
	//读取kpmfile
	kpmfilep, err := NewKpmFileP(pwd)
	if err != nil {
		return err
	}
	direct := kpmfilep.kpmfile.Direct
	directMap := make(map[string]Require, 16)
	for i := 0; i < len(direct); i++ {
		if direct[i].Alias == "" {
			if direct[i].Type != "git" {
				directMap[direct[i].Name] = direct[i]
			}
		} else {
			directMap[direct[i].Alias] = direct[i]
		}

	}
	//间接依赖，添加原则，唯一版本即可
	indirect := kpmfilep.kpmfile.Indirect
	indirectMap := make(map[string]Require, 16)
	for i := 0; i < len(indirect); i++ {
		indirectMap[indirect[i].Type+"|"+indirect[i].Name+"|"+indirect[i].GitAddress+"|"+indirect[i].Version+"|"+indirect[i].GitCommit] = indirect[i]
	}
	for i := 0; i < len(pkgvs); i++ {

		r := &Require{}
		err := r.NewRequireFromPkgString(pkgvs[i], flag_git)
		if err != nil {
			return err
		}

		err = r.Get(KPM_ROOT, KPM_SERVER_ADDR)
		if err != nil {
			return err
		}
		//检查命名是否冲突，先读，后写
		if r.Alias == "" {
			if r.Type != "git" {
				//仓库包
				_, stat := directMap[r.Name]
				if stat {
					//冲突
					println("Naming conflicts")
					continue
				}
				directMap[r.Name] = *r
			}
		} else {
			_, stat := directMap[r.Alias]
			if stat {
				//冲突
				println("Naming conflicts")
				continue
			}
			directMap[r.Alias] = *r
			file, err := os.ReadFile(r.KpmFileLocalPath(KPM_ROOT, KPM_SERVER_ADDR_PATH))
			if err == nil {
				//解析
				kpmfile := KpmFile{}
				err = json.Unmarshal(file, &kpmfile)
				if err != nil {
					return err
				}
				//检查kcl版本，如果高于当前版本则拒绝
				//工作版本
				ver := &Version{}
				err = ver.NewFromString(kpmfilep.kpmfile.KclvmMinVersion)
				if err != nil {
					return err
				}
				//当前解析依赖的版本
				nowver := &Version{}
				err = ver.NewFromString(kpmfile.KclvmMinVersion)
				if err != nil {
					return err
				}
				if ver.Cmp(*nowver) == -1 {
					//println("The current pending load dependency aabc needs to be greater than "+kpmfile.KclvmMinVersion+" version of KclvmMinVersion")
					return errors.New("The current pending load dependency aabc needs to be greater than " + kpmfile.KclvmMinVersion + " version of KclvmMinVersion")
				}
				//遍历得到直接依赖和间接依赖
				for j := 0; j < len(kpmfile.Direct); j++ {
					tmp := kpmfile.Direct[j]
					indirectMap[tmp.Type+"|"+tmp.Name+"|"+tmp.GitAddress+"|"+tmp.Version+"|"+tmp.GitCommit] = tmp
				}
				for j := 0; j < len(kpmfile.Indirect); j++ {
					tmp := kpmfile.Indirect[j]
					indirectMap[tmp.Type+"|"+tmp.Name+"|"+tmp.GitAddress+"|"+tmp.Version+"|"+tmp.GitCommit] = tmp
				}

			}
			//没有文件，不需要解析依赖
			//return err

		}
		err = r.LinkToExternal(KPM_ROOT, KPM_SERVER_ADDR_PATH, pwd)
		if err != nil {
			println(err.Error())
			return err
		}
	}
	//回填依赖数据并保存
	kpmfilep.kpmfile.Direct = kpmfilep.kpmfile.Direct[:0]
	for _, v := range directMap {
		kpmfilep.kpmfile.Direct = append(kpmfilep.kpmfile.Direct, v)
	}
	kpmfilep.kpmfile.Indirect = kpmfilep.kpmfile.Indirect[:0]
	for _, v := range indirectMap {
		kpmfilep.kpmfile.Indirect = append(kpmfilep.kpmfile.Indirect, v)
	}
	if debuglog {
		fmt.Println("directMap", directMap)
	}

	err = kpmfilep.Save()
	if err != nil {
		return err
	}
	return nil
}

// CliDel 移除链接,删除直接依赖的包信息,别名
func CliDel(args ...string) error {
	kpmfilep, err := NewKpmFileP(pwd)
	if err != nil {
		return err
	}
	direct := kpmfilep.kpmfile.Direct
	directMap := make(map[string]Require, 16)
	for i := 0; i < len(direct); i++ {
		if direct[i].Alias == "" {
			if direct[i].Type != "git" {
				directMap[direct[i].Name] = direct[i]
			}
		} else {
			directMap[direct[i].Alias] = direct[i]
		}

	}
	for i := 0; i < len(args); i++ {
		t, stat := directMap[args[i]]
		if !stat {
			//
			//println("del  dependencies", args[i], " fail,it does not exist in kpmfile")
			return errors.New("del  dependencies " + args[i] + " fail,it does not exist in kpmfile")
		}
		name := t.Alias
		if t.Alias == "" {
			name = t.Name
		}
		err = os.Remove(pwd + Separator + ExternalDependencies + Separator + name)
		if err != nil {
			println("del  dependencies", name, " fail")
			return err
		}
		delete(directMap, args[i])
		println("del  dependencies", name, " success")
	}
	kpmfilep.kpmfile.Direct = kpmfilep.kpmfile.Direct[:0]
	for _, v := range directMap {
		kpmfilep.kpmfile.Direct = append(kpmfilep.kpmfile.Direct, v)
	}

	err = kpmfilep.Save()
	if err != nil {
		return err
	}
	return nil
}

func CliInit(pkg string) error {
	kpmfp := &KpmFileP{
		Path: pwd + Separator + "kpm.json",
		kpmfile: &KpmFile{
			PackageName:     pkg,
			KclvmMinVersion: KclvmMinVersion,
		},
	}
	err := kpmfp.Create()
	if err != nil {
		return errors.New("Create kpm.json fail!Because " + err.Error())
	}

	println("Create kpm.json success!")
	_, err = os.Stat(pwd + Separator + "kcl.mod")
	if err == nil {
		return nil
	}
	//文件不存在,所以创建
	err = os.WriteFile(pwd+Separator+"kcl.mod", []byte(DefaultKclModContent+`"`+KclvmMinVersion+`"`), 0777)
	if err != nil {
		return err
	}
	return nil
}
