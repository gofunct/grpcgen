package project

import (
	"github.com/gofunct/grpcgen/logging"
	"os"
	"path/filepath"
	"strings"
)

type Project struct {
	AbsPath string
	CmdPath string
	SrcPath string
	Name    string
}

// NewProject returns Project with specified project Name.
func NewProject(projectName string) *Project {
	if projectName == "" {
		logging.Exit("can't create project with blank Name")
	}

	p := new(Project)
	p.Name = projectName

	// 1. Find already created protect.
	p.AbsPath = FindPackage(projectName)

	// 2. If there are no created project with this path, and user is in GOPATH,
	// then use GOPATH/src/projectName.
	if p.AbsPath == "" {
		wd, err := os.Getwd()
		logging.IfErr("failed to get working directory", err)
		for _, SrcPath := range srcPaths {
			goPath := filepath.Dir(SrcPath)
			if FilePathHasPrefix(wd, goPath) {
				p.AbsPath = filepath.Join(SrcPath, projectName)
				break
			}
		}
	}

	// 3. If user is not in GOPATH, then use (first GOPATH)/src/projectName.
	if p.AbsPath == "" {
		p.AbsPath = filepath.Join(srcPaths[0], projectName)
	}

	return p
}

// NewProjectFromPath returns Project with specified absolute path to
// package.
func NewProjectFromPath(AbsPath string) *Project {
	if AbsPath == "" {
		logging.Exit("can't create project: AbsPath can't be blank")
	}
	if !filepath.IsAbs(AbsPath) {
		logging.Exit("can't create project: AbsPath is not absolute")
	}

	// If AbsPath is symlink, use its destination.
	fi, err := os.Lstat(AbsPath)
	logging.IfErr("can't read path info: ", err)

	if fi.Mode()&os.ModeSymlink != 0 {
		path, err := os.Readlink(AbsPath)
		logging.IfErr("can't read the destination of symlink: ", err)
		AbsPath = path
	}

	p := new(Project)
	p.AbsPath = strings.TrimSuffix(AbsPath, FindCmdDir(AbsPath))
	p.Name = filepath.ToSlash(TrimScrcPath(p.AbsPath, p.GetSource()))
	return p
}

// Name returns the Name of project, e.g. "github.com/spf13/cobra"
func (p *Project) GetName() string {
	return p.Name
}

// CmdPath returns absolute path to directory, where all commands are located.
func (p *Project) GetCmd() string {
	if p.AbsPath == "" {
		return ""
	}
	if p.CmdPath == "" {
		p.CmdPath = filepath.Join(p.AbsPath, FindCmdDir(p.AbsPath))
	}
	return p.CmdPath
}

func InitializeProject(p *Project) {
	if !PathExists(p.GetAbsPath()) { // If path doesn't yet exist, create it
		err := os.MkdirAll(p.GetAbsPath(), os.ModePerm)
		logging.IfErr("failed to make directories", err)

	} else if !EmptyPath(p.GetAbsPath()) { // If path exists and is not empty don't use it
		logging.Exit("Gen will not create a new project in a non empty directory: " + p.GetAbsPath())
	}

	// We have a directory and it's empty. Time to initialize it.
	p.CreateMainFile()
	p.CreateDockerfile()
	p.CreateMakeFile()
	p.CreateProtofile()
	p.CreateRootCmdFile()
}

// AbsPath returns absolute path of project.
func (p Project) GetAbsPath() string {
	return p.AbsPath
}

// AbsPath returns absolute path of project.
func (p Project) Absolute() string {
	return p.AbsPath
}

// SrcPath returns absolute path to $GOPATH/src where project is located.
func (p *Project) GetSource() string {
	if p.SrcPath != "" {
		return p.SrcPath
	}
	if p.AbsPath == "" {
		p.SrcPath = srcPaths[0]
		return p.SrcPath
	}

	for _, SrcPath := range srcPaths {
		if FilePathHasPrefix(p.AbsPath, SrcPath) {
			p.SrcPath = SrcPath
			break
		}
	}

	return p.SrcPath
}
