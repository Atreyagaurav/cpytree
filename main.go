package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Tree struct {
	Parent   *Tree
	Rank     int
	ChildNum int
	Value    string
	Children map[string]*Tree
}

func (t *Tree) AddChild(cstr string) *Tree {
	ch, ok := t.Children[cstr]
	if ok {
		return ch
	} else {
		c := &Tree{Value: cstr, Children: make(map[string]*Tree)}
		c.Parent = t
		t.Children[cstr] = c
		t.ChildNum += 1
		c.Rank = t.Rank + 1
		return c
	}
}

// func (t *Tree) ConnectChild(c *Tree) *Tree {
// 	t.ChildNum += 1
// 	t.Children = append(t.Children, c)
// 	c.Parent = t
// 	c.Rank = t.Rank + 1
// 	return c
// }

func (t *Tree) AddChildFromString(path string) *Tree {
	vals := strings.Split(path, "/")
	c := t
	for _, p := range vals {
		c = c.AddChild(p)
	}
	return c
}

func (t *Tree) Construct() {
	for _, c := range t.Children {
		c.Construct()
	}
	if t.ChildNum == 0 {
		path := t.GetFull("/")
		file, err := os.OpenFile("makedir.sh", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			log.Print("Error in reading file:", "makedir.sh")
		}
		file.WriteString(fmt.Sprintf("mkdir -p %s\n", path))
		file.Close()
		// os.MkdirAll(path, 0700)
	}
}

func (t *Tree) GetFull(s string) string {
	path := t.Value
	if t.Parent == nil {
		// cw, _ := os.Getwd()
		// return cw + s + path
		return path
	} else {
		return fmt.Sprintf("%s%s%s", t.Parent.GetFull(s), s, t.Value)
	}
}

func (t *Tree) Show() {
	for i := 0; i < t.Rank; i++ {
		fmt.Print("  ")
	}
	fmt.Printf("%s\n", t.Value)
	for _, c := range t.Children {
		c.Show()
	}
}

func main() {
	args := os.Args[1:]
	thispath, _ := os.Getwd()
	// MakeDirectoryTree(args)
	var pathList []string
	t := &Tree{Value: ".", Children: make(map[string]*Tree)}
	if len(args) == 1 {
		filepath.Walk(args[0], func(path string, info os.FileInfo, err error) error {
			if filepath.IsAbs(path) {
				path, _ = filepath.Rel(thispath, path)
			}
			fi, _ := os.Stat(path)
			if fi.IsDir() && path[0:1] != "." {
				pathList = append(pathList, path)
				t.AddChildFromString(path)
			}
			return nil
		})
	} else {
		for _, v := range args {
			t.AddChildFromString(v)
		}
	}
	// t.Show()
	os.Remove("makedir.sh")
	t.Construct()
	return
}
