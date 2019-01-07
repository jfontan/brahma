package brahma

import (
	"os"
	"path/filepath"

	sivafs "gopkg.in/src-d/go-billy-siva.v4"
	"gopkg.in/src-d/go-billy.v4/osfs"
	"gopkg.in/src-d/go-billy.v4/util"
	errors "gopkg.in/src-d/go-errors.v1"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing/cache"
	"gopkg.in/src-d/go-git.v4/storage/filesystem"
)

const (
	FetchRefSpec = config.RefSpec("refs/*:refs/*")
	FetchHEAD    = config.RefSpec("HEAD:refs/heads/HEAD")
)

var (
	ErrFileAlreadyExists = errors.NewKind("file already exists: %s")
	ErrDownloadingRepo   = errors.NewKind("error downloading repo")
)

func Download(url, path string) error {
	sivaFile, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	_, err = os.Stat(sivaFile)
	if err == nil {
		return ErrFileAlreadyExists.New(sivaFile)
	}

	rootFS := osfs.New("/")

	brahmaFS := osfs.New("/tmp/brahma")
	tmpDir, err := util.TempDir(brahmaFS, ".", "download")
	if err != nil {
		return ErrDownloadingRepo.Wrap(err)
	}

	defer util.RemoveAll(brahmaFS, tmpDir)

	tmpFS, err := brahmaFS.Chroot(tmpDir)
	if err != nil {
		return ErrDownloadingRepo.Wrap(err)
	}

	siva, err := sivafs.NewFilesystem(rootFS, sivaFile, tmpFS)
	if err != nil {
		return ErrDownloadingRepo.Wrap(err)
	}

	storage := filesystem.NewStorage(siva, cache.NewObjectLRUDefault())
	repo, err := git.Init(storage, nil)
	if err != nil {
		return ErrDownloadingRepo.Wrap(err)
	}

	remote, err := repo.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{url},
	})
	if err != nil {
		return ErrDownloadingRepo.Wrap(err)
	}

	o := &git.FetchOptions{
		RefSpecs: []config.RefSpec{FetchRefSpec, FetchHEAD},
		Force:    true,
	}
	err = remote.Fetch(o)
	if err != nil {
		return ErrDownloadingRepo.Wrap(err)
	}

	err = siva.Sync()
	if err != nil {
		rootFS.Remove(sivaFile)
		return ErrDownloadingRepo.Wrap(err)
	}

	return nil
}
