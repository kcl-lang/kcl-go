package main

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/orangebees/go-oneutils/PathHandle"
	"os"
	"strings"
)

func main() {

	err := PathHandle.RunInTempDir(func(tmppath string) error {
		println(tmppath)
		r, err := git.PlainClone(tmppath, false, &git.CloneOptions{
			URL:      "https://" + "github.com/go-git/go-git",
			Progress: os.Stdout,
		})
		iter, err := r.Tags()
		if err != nil {
			return err
		}
		var commitHash = make([]byte, 40)
		commitHash = commitHash[:0]
		var commitTag = make([]byte, 8)
		commitTag = commitTag[:0]
		err = iter.ForEach(func(ref *plumbing.Reference) error {
			commitHash = append(commitHash[:0], ref.Hash().String()...)
			commitTag = append(commitTag[:0], strings.TrimPrefix(string(ref.Name()), "refs/tags/")...)
			return nil
		})
		println(string(commitHash))
		println(string(commitTag))
		//time.Sleep(time.Second * 100)
		return nil
	})
	if err != nil {
		return
	}
	//r, err := git.PlainClone("C:\\aaaaa", false, &git.CloneOptions{
	//	URL:      "https://github.com/go-git/go-git",
	//	Progress: os.Stdout,
	//})
	//r, err := git.PlainOpen("C:\\aaaaa")
	//if err != nil {
	//
	//}
	//i := 0
	//cIter, _ := r.Log(&git.LogOptions{From: plumbing.Hash{}})
	//err = cIter.ForEach(func(c *object.Commit) error {
	//	if i == 0 {
	//
	//		println(c.ID().String())
	//		println(c.Committer.When.Unix())
	//	}
	//	if c.ID().String() == "da810275bf682d29a530ed819aff175f47bd7634" {
	//		fmt.Println(c)
	//	}
	//	i++
	//	return nil
	//})

	//commitObject, err := r.CommitObject(plumbing.NewHash("452df976faca7193b275bd31a0d76027a2a9df5c"))
	//if err != nil {
	//	println(err.Error())
	//	return
	//}
	//fmt.Println(commitObject)
	//iter, err := r.Tags()
	//if err != nil {
	//	return
	//}
	//var commitHash = make([]byte, 40)
	//commitHash = commitHash[:0]
	//var commitTag = make([]byte, 8)
	//commitTag = commitTag[:0]
	//err = iter.ForEach(func(ref *plumbing.Reference) error {
	//	commitHash = append(commitHash[:0], ref.Hash().String()...)
	//	commitTag = append(commitTag[:0], strings.TrimPrefix(string(ref.Name()), "refs/tags/")...)
	//	return nil
	//})
	//println(string(commitHash))
	//println(string(commitTag))
	//if err != nil {
	//	return
	//}
	//先获取最新tag，如果没有tag，则直接使用hash
	//err := kpm.CLI(os.Args...)
	//if err != nil {
	//	println(err.Error())
	//}

}
