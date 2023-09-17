package filesystem

import (
	"context"
	"github.com/pkg/errors"
	"golang.org/x/net/webdav"
	"os"
)

type Disk struct {
	root string
}

func New(root string) webdav.FileSystem {
	return &Disk{
		root: root,
	}
}

func (f Disk) Mkdir(_ context.Context, name string, perm os.FileMode) error {
	err := os.Chdir(f.root)
	if err != nil {
		return errors.Wrap(err, "failed to change directory")
	}

	return errors.Wrap(os.Mkdir(name, perm), "failed to make directory")
}

func (f Disk) OpenFile(_ context.Context, name string, flag int, perm os.FileMode) (webdav.File, error) {
	err := os.Chdir(f.root)
	if err != nil {
		return nil, errors.Wrap(err, "failed to change directory")
	}

	file, err := os.OpenFile(name, flag, perm)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open file")
	}

	return file, nil
}

func (f Disk) RemoveAll(_ context.Context, name string) error {
	err := os.Chdir(f.root)
	if err != nil {
		return errors.Wrap(err, "failed to change directory")
	}

	return errors.Wrap(os.RemoveAll(name), "failed to remove all")
}

func (f Disk) Rename(_ context.Context, oldName, newName string) error {
	err := os.Chdir(f.root)
	if err != nil {
		return errors.Wrap(err, "failed to change directory")
	}

	return errors.Wrap(os.Rename(oldName, newName), "failed to rename")
}

func (f Disk) Stat(_ context.Context, name string) (os.FileInfo, error) {
	err := os.Chdir(f.root)
	if err != nil {
		return nil, errors.Wrap(err, "failed to change directory")
	}

	stat, err := os.Stat(name)
	return stat, errors.Wrap(err, "failed to get stat")
}
