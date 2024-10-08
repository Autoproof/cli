package project

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const (
	// AutoproofHomeDirName is a name of the directory where autoproof stores its data.
	AutoproofHomeDirName = ".autoproof"
)

// Project is a representation of the autoproof project.
type Project struct {
	name string
	path string

	config *Config
}

type SnapshotItem struct {
	Filename string
	Hash     string
}

type Snapshot struct {
	Items []SnapshotItem
}

// FromPath returns a Project from the given path.
// Path must be a path to the project directory (i.e. a directory where .autoproof directory is located).
func FromPath(projectPath string) (*Project, error) {
	projectPath, err := filepath.Abs(projectPath)
	if err != nil {
		return nil, err
	}

	fi, err := os.Lstat(filepath.Join(projectPath, AutoproofHomeDirName))
	switch {
	case os.IsNotExist(err):
		return nil, errors.New("the path to the configuration must be exists")
	case err != nil:
		return nil, err
	}

	if !fi.IsDir() {
		return nil, fmt.Errorf("%s should be a directory", AutoproofHomeDirName)
	}
	p := &Project{
		path: projectPath,
	}

	if err := p.readConfig(); err != nil {
		return nil, fmt.Errorf("read project config: %w", err)
	}

	p.name = p.Config().Get("projectName").(string)
	return p, nil
}

func New(name, path string) (*Project, error) {
	p := &Project{
		name: name,
		path: path,
	}

	if err := os.MkdirAll(filepath.Join(path, AutoproofHomeDirName), 0777); err != nil {
		return nil, fmt.Errorf("create project dir: %w", err)
	}

	if err := p.readConfig(); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	p.Config().Set("projectName", name)
	return p, nil
}

func (p *Project) Name() string {
	return p.name
}
