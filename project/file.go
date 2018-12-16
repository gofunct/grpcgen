package project

import (
	"bytes"
	"fmt"
	"go/format"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/xlab/treeprint"
)

type File struct {
	Name     string
	AbsPath  string
	Template string
}

type Folder struct {
	Name    string
	AbsPath string

	// Unexported so you can't set them without methods
	files   []File
	folders []*Folder
}

func (f *Folder) AddFolder(name string) *Folder {
	newF := &Folder{
		Name:    name,
		AbsPath: filepath.Join(f.AbsPath, name),
	}
	f.folders = append(f.folders, newF)
	return newF
}

func (f *Folder) AddFile(name, tmpl string) {
	f.files = append(f.files, File{
		Name:     name,
		Template: tmpl,
		AbsPath:  filepath.Join(f.AbsPath, name),
	})
}

func (f *Folder) RenderTemplate(templatePath string, p Project) error {
	for _, v := range f.files {
		t, err := template.ParseFiles(filepath.Join(templatePath, v.Template))
		if err != nil {
			return err
		}

		File, err := os.Create(v.AbsPath)
		if err != nil {
			return err
		}

		defer File.Close()

		if strings.Contains(v.AbsPath, ".go") {
			var out bytes.Buffer
			err = t.Execute(&out, p)
			if err != nil {
				log.Printf("Could not process template %s\n", v)
				return err
			}

			b, err := format.Source(out.Bytes())
			if err != nil {
				fmt.Print(string(out.Bytes()))
				log.Printf("\nCould not format Go File %s\n", v)
				return err
			}

			_, err = File.Write(b)
			if err != nil {
				return err
			}
		} else {
			err = t.Execute(File, p)
			if err != nil {
				return err
			}
		}
	}

	for _, v := range f.folders {
		err := os.Mkdir(v.AbsPath, os.ModePerm)
		if err != nil {
			return err
		}

		err = v.RenderTemplate(templatePath, p)
		if err != nil {
			return err
		}
	}

	return nil
}

func (f Folder) PrintTree() {
	t := f.GetTree(true, treeprint.New())
	fmt.Println(t.String())
}

func (f Folder) GetTree(root bool, tree treeprint.Tree) treeprint.Tree {
	if !root {
		tree = tree.AddBranch(f.Name)
	}

	for _, v := range f.folders {
		v.GetTree(false, tree)
	}

	for _, v := range f.files {
		tree.AddNode(v.Name)
	}

	return tree
}
