package scripts

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	kcl_plugin "kusionstack.io/kcl-plugin"
)

type KclvmAssetHelper struct {
	Triple  KclvmTripleType
	Version KclvmVersionType
}

func NewKclvmAssetHelper(kclvmTriple KclvmTripleType, kclvmVersion KclvmVersionType) *KclvmAssetHelper {
	if kclvmTriple == "" {
		panic(fmt.Errorf("kclvm triple missing"))
	}
	if kclvmVersion == "" {
		panic(fmt.Errorf("kclvm version missing"))
	}

	return &KclvmAssetHelper{
		Version: kclvmVersion,
		Triple:  kclvmTriple,
	}
}

func (p *KclvmAssetHelper) GetFilename() string {
	if p.Triple == KclvmTripleType_windows {
		return fmt.Sprintf("kclvm-%s-%s.zip", p.Version, p.Triple)
	} else {
		return fmt.Sprintf("kclvm-%s-%s.tar.gz", p.Version, p.Triple)
	}
}

func (p *KclvmAssetHelper) GetFileMd5um() string {
	kclvmFilename := p.GetFilename()
	return KclvmMd5sum[kclvmFilename]
}

func (p *KclvmAssetHelper) GetDownloadUrl(baseUrl string) string {
	baseUrl = strings.TrimRight(baseUrl, "/")
	return fmt.Sprintf("%s/%s/%s", baseUrl, p.Version, p.GetFilename())
}

func (p *KclvmAssetHelper) DownloadFile(localFilename string) error {
	md5sum := p.GetFileMd5um()
	if md5sum == "" {
		return fmt.Errorf("%s: not found, md5sum missing", p.GetFilename())
	}
	if MD5File(localFilename) == md5sum {
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // cancel when we are finished consuming integers

	var urls = []string{
		p.GetDownloadUrl(KclvmDownloadUrlBase),
	}

	for _, s := range KclvmDownloadUrlBase_mirrors {
		mirrorBase := strings.TrimSpace(s)
		if mirrorBase != "" {
			urls = append(urls, p.GetDownloadUrl(mirrorBase))
		}
	}

	var errs = make(chan error, len(urls))
	var okfiles = make(chan string, len(urls))
	var wg sync.WaitGroup

	wg.Add(len(urls))
	for i, s := range urls {
		go func(id int, url, localFilename string) {
			defer wg.Done()
			tmpname := fmt.Sprintf("%s.%d", localFilename, id)
			err := HttpGetFile(ctx, url, tmpname, true)
			if err != nil {
				errs <- err
				return
			}
			if got := MD5File(tmpname); got != md5sum {
				errs <- fmt.Errorf("md5 mismatch: expect=%v, got=%v, url=%s", md5sum, got, url)
				return
			}

			// OK
			okfiles <- tmpname
			cancel()
		}(i, s, localFilename)
	}
	wg.Wait()

	if len(okfiles) > 0 {
		tmpname := <-okfiles
		os.Rename(tmpname, localFilename)

		for id := range urls {
			tmpname := fmt.Sprintf("%s.%d", localFilename, id)
			os.Remove(tmpname)
		}

		if got := MD5File(localFilename); got != md5sum {
			return fmt.Errorf("md5 mismatch: expect=%v, got=%v, local=%s", md5sum, got, localFilename)
		}

		return nil
	}

	return <-errs
}

func (p *KclvmAssetHelper) Install(kclvmRoot string) (err error) {
	md5sumFile := filepath.Join(kclvmRoot, "md5sum.txt")
	if FileExists(md5sumFile) {
		return nil
	}

	var localFilename = "zz_download-" + p.GetFilename()
	defer func() {
		if err == nil {
			os.Remove(localFilename)
		}
	}()

	if err := p.DownloadFile(localFilename); err != nil {
		return err
	}

	if strings.HasSuffix(localFilename, ".zip") {
		if err := Unzip(localFilename, kclvmRoot); err != nil {
			return err
		}
	} else {
		if err := UnTarGz(localFilename, "kclvm", kclvmRoot); err != nil {
			return err
		}
	}

	// chmod +x
	for _, s := range kclvm_bin_exe_list {
		basename := s
		if p.Triple == KclvmTripleType_windows {
			basename += ".exe"
		}
		filepath := filepath.Join(kclvmRoot, "bin", basename)
		err := os.Chmod(filepath, 0777)

		fmt.Println("chmod", filepath, err)
	}

	// write md5sum
	if s := filepath.Join(kclvmRoot, "md5sum.txt"); !FileExists(s) {
		txt := fmt.Sprintf("%s *%s\n", p.GetFileMd5um(), p.GetFilename())
		if err := ioutil.WriteFile(s, []byte(txt), 0666); err != nil {
			return err
		}
	}

	// write VERSION
	if s := filepath.Join(kclvmRoot, "VERSION"); !FileExists(s) {
		if err := ioutil.WriteFile(s, []byte(DefaultKclvmVersion), 0666); err != nil {
			return err
		}
	}

	kclvmPluginsPath := p.getPluginPath(kclvmRoot)
	if err := kcl_plugin.InstallPlugins(kclvmPluginsPath); err != nil {
		return err
	}

	return nil
}

func (p *KclvmAssetHelper) getPluginPath(kclvmRoot string) string {
	if p.Triple == KclvmTripleType_windows {
		kclvmPluginPath := filepath.Join(kclvmRoot, "bin", "plugins")
		return kclvmPluginPath
	}
	kclvmPluginPath := filepath.Join(kclvmRoot, "plugins")
	return kclvmPluginPath
}

var kclvm_bin_exe_list = []string{
	"kclvm",
	"kclvm_cli",

	"kcl",
	"kcl-plugin",
	"kcl-doc",
	"kcl-test",
	"kcl-lint",
	"kcl-fmt",
	"kcl-vet",
}
