package viperizer

import (
	"os"
	"path/filepath"
	"runtime"
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
		er("can't create project with blank Name")
	}

	p := new(Project)
	p.Name = projectName

	// 1. Find already created protect.
	p.AbsPath = findPackage(projectName)

	// 2. If there are no created project with this path, and user is in GOPATH,
	// then use GOPATH/src/projectName.
	if p.AbsPath == "" {
		wd, err := os.Getwd()
		if err != nil {
			er(err)
		}
		for _, SrcPath := range srcPaths {
			goPath := filepath.Dir(SrcPath)
			if filepathHasPrefix(wd, goPath) {
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

// findPackage returns full path to existing go package in GOPATHs.
func findPackage(packageName string) string {
	if packageName == "" {
		return ""
	}

	for _, SrcPath := range srcPaths {
		packagePath := filepath.Join(SrcPath, packageName)
		if exists(packagePath) {
			return packagePath
		}
	}

	return ""
}

// NewProjectFromPath returns Project with specified absolute path to
// package.
func NewProjectFromPath(AbsPath string) *Project {
	if AbsPath == "" {
		er("can't create project: AbsPath can't be blank")
	}
	if !filepath.IsAbs(AbsPath) {
		er("can't create project: AbsPath is not absolute")
	}

	// If AbsPath is symlink, use its destination.
	fi, err := os.Lstat(AbsPath)
	if err != nil {
		er("can't read path info: " + err.Error())
	}
	if fi.Mode()&os.ModeSymlink != 0 {
		path, err := os.Readlink(AbsPath)
		if err != nil {
			er("can't read the destination of symlink: " + err.Error())
		}
		AbsPath = path
	}

	p := new(Project)
	p.AbsPath = strings.TrimSuffix(AbsPath, findCmdDir(AbsPath))
	p.Name = filepath.ToSlash(trimSrcPath(p.AbsPath, p.GetSource()))
	return p
}

// trimSrcPath trims at the beginning of AbsPath the SrcPath.
func trimSrcPath(AbsPath, SrcPath string) string {
	relPath, err := filepath.Rel(SrcPath, AbsPath)
	if err != nil {
		er(err)
	}
	return relPath
}

// Name returns the Name of project, e.g. "github.com/spf13/cobra"
func (p Project) GetName() string {
	return p.Name
}

// CmdPath returns absolute path to directory, where all commands are located.
func (p *Project) GetCmd() string {
	if p.AbsPath == "" {
		return ""
	}
	if p.CmdPath == "" {
		p.CmdPath = filepath.Join(p.AbsPath, findCmdDir(p.AbsPath))
	}
	return p.CmdPath
}

// findCmdDir checks if base of AbsPath is cmd dir and returns it or
// looks for existing cmd dir in AbsPath.
func findCmdDir(AbsPath string) string {
	if !exists(AbsPath) || isEmpty(AbsPath) {
		return "cmd"
	}

	if isCmdDir(AbsPath) {
		return filepath.Base(AbsPath)
	}

	files, _ := filepath.Glob(filepath.Join(AbsPath, "c*"))
	for _, file := range files {
		if isCmdDir(file) {
			return filepath.Base(file)
		}
	}

	return "cmd"
}

// isCmdDir checks if base of Name is one of cmdDir.
func isCmdDir(Name string) bool {
	Name = filepath.Base(Name)
	for _, cmdDir := range []string{"cmd", "cmds", "command", "commands"} {
		if Name == cmdDir {
			return true
		}
	}
	return false
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
		if filepathHasPrefix(p.AbsPath, SrcPath) {
			p.SrcPath = SrcPath
			break
		}
	}

	return p.SrcPath
}

func filepathHasPrefix(path string, prefix string) bool {
	if len(path) <= len(prefix) {
		return false
	}
	if runtime.GOOS == "windows" {
		// Paths in windows are case-insensitive.
		return strings.EqualFold(path[0:len(prefix)], prefix)
	}
	return path[0:len(prefix)] == prefix

}

// AbsPath returns absolute path of project.
func (p Project) GetAbsPath() string {
	return p.AbsPath
}
