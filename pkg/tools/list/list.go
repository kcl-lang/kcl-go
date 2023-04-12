// Copyright 2023 The KCL Authors. All rights reserved.

/*
Package list extracts information by parsing KCL source code and return it as list. For now, the supported information to be listed will be UpStream/DownStream dependency files.
It can also be schemas, schema attributes and so on. Supporting on listing these kinds of information is in the plan.
*/
package list

// ListDepFiles return the depend files from the given path. It will scan and parse the kusion applications within the workdir,
// then list depend files of the applications.
func ListDepFiles(workDir string, opt *Option) (files []string, err error) {
	if opt == nil {
		opt = &Option{}
	}

	pkgroot, pkgpath, err := FindPkgInfo(workDir)
	if err != nil {
		return nil, err
	}

	depParser := NewSingleAppDepParser(pkgroot, *opt)

	appFiles, err := depParser.GetAppFiles(pkgpath, opt.FlagAll)
	if err != nil {
		return nil, err
	}
	for _, s := range appFiles {
		if opt.UseAbsPath {
			files = append(files, pkgroot+"/"+s)
		} else {
			files = append(files, s)
		}
	}

	return files, nil
}

// ListUpStreamFiles returns a list of the UpStream dependent packages/files from the given path list.
//
// Usage Caution
//
// The implementation of this API is based on reading files in opt.Files so the time-consuming is positively related to the number of files.
// Do not call the API with high frequency and please ensure at least 10 seconds interval when calling.
//
// Terminology
//
// The word "UpStream" means the dependent direction between two files/packages. One file/package (named f) depends on its UpStream files/packages
// if there exist some directly or indirectly import statements that finally import the UpStream files/packages to f.
//
// Parameters
//
// The "workDir" must be a valid KCL program root directory, otherwise a "pkgroot: not found" error will be produced.
// The param opt.Files specifies the scope of the KCL files to be analyzed and list UpStream files/packages on.
// The param opt.UpStreams can be set nil or empty. It will not be used.
//
// The API will return nil if opt or opt.Files is nil.
//
// Example
//
// For instance a KCL program that comes with three files: main.k and base/a.k, base/b.k,
// and the file main.k contains an import statement that imports base/b.k to it, while the file base/b.k imports base/a.k:
//
// 	demo (KCL program root)
// 	├── base
// 	│   ├── a.k
// 	│   └── b.k         # import .a
// 	└── main.k          # import base.b
//
// Then the UpStream files/packages of the file main.k will be: base/a.k and base/b.k
// If the import statement in main.k changes to:
//
// 	import base
//
// To list UpStream files/packages of the file main.k, the function call will be:
//
//  ListUpStreamFiles("demo", &DepOptions{Files:[]string{"main.k"})
//
// Then its UpStream files/packages will be: base, base/a.k and base/b.k
func ListUpStreamFiles(workDir string, opt *DepOptions) (deps []string, err error) {
	if opt == nil || len(opt.Files) == 0 {
		return nil, nil
	}
	pkgroot, _, err := FindPkgInfo(workDir)
	if err != nil {
		return nil, err
	}
	depParser, err := newImportDepParser(pkgroot, *opt)
	if err != nil {
		return nil, err
	}
	return depParser.upstreamFiles(), nil
}

// ListDownStreamFiles returns a list of DownStream dependent packages/files from the given changed path list.
// A typical use is to list all the DownStream files when some files changed(added/modified/deleted) in a KCL configuration repository so that
// certain test cases on those files, instead of all the test cases will need to be rerun to save integration time.
//
// Usage Caution
//
// The implementation of this API is based on reading files in opt.Files so the time-consuming is positively related to the number of files.
// Do not call the API with high frequency and please ensure at least 10 seconds interval when calling.
//
// Terminology
//
// The word "DownStream" means the dependent direction between two files/packages. One file/package is dependent by its DownStream files/packages
// if there exist some directly or indirectly import statements in those files/packages that finally import f to them.
//
// Parameters
//
// The "workDir" must be a valid KCL program root directory, otherwise a "pkgroot: not found" error will be produced.
//
// The param opt.Files specifies the scope of the KCL files to be analyzed.
// Thus only files/packages that have directly or indirectly UpStream/DownStream relations with those files will appears in the result.
//
// The param opt.UpStreams specifies the KCL files to list DownStream files/packages on.
//
// The API will return nil if either opt/opt.Files/opt.UpStreams is nil.
//
// Example
//
// For instance, a KCL program that comes with three files: main.k and base/a.k, base/b.k,
// and the file main.k contains an import statement that imports base/b.k to it, while the file base/b.k imports base/a.k:
//
// 	demo (KCL program root)
// 	├── base
// 	│   ├── a.k
// 	│   └── b.k         # import .a
// 	└── main.k          # import base.b
//
// To list DownStream files/packages of the file base/a.k, the function call will be:
//
//  ListDownStreamFiles("demo", &DepOptions{Files:[]string{"main.k"}, UpStreams:[]string{"base/a.k"})
//
// Then the DownStream files/packages of the file base/a.k will be: base/b.k, base and main.k
// If the import statement in main.k changes to:
//
// 	import base
//
// Then its DownStream files will be: base, base/a.k and base/b.k
// And if the import statement in main.k changes to:
//
// 	import base
// And the DownStream files/packages of the file base/a.k stays the same
//
func ListDownStreamFiles(workDir string, opt *DepOptions) ([]string, error) {
	if opt == nil || len(opt.Files) == 0 || len(opt.UpStreams) == 0 {
		return nil, nil
	}
	pkgroot, _, err := FindPkgInfo(workDir)
	if err != nil {
		return nil, err
	}
	depParser, err := newImportDepParser(pkgroot, *opt)
	if err != nil {
		return nil, err
	}

	return depParser.downStreamFiles(), nil
}
