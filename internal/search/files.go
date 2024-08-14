package search

import (
	"io"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
)

// ListFiles returns a list of available files for searching
//
// if the given root is a git repository, the returned list is fetched from the git index (`git ls-files`)
// if the given root is not a git repository, the returned list is all files excluding `node_modules`
func ListFiles(root string) ([]string, error) {
	repo, err := git.PlainOpen(root)
	if err != nil {
		return listOsFiles(root)
	}
	return listGitFiles(repo)
}

func listGitFiles(repo *git.Repository) ([]string, error) {
	ref, err := repo.Head()
	if err != nil {
		return nil, err
	}

	commit, err := repo.CommitObject(ref.Hash())
	if err != nil {
		return nil, err
	}

	files, err := commit.Files()
	if err != nil {
		return nil, err
	}

	var gitFiles []string
	for {
		file, err := files.Next()
		if err != nil && err == io.EOF {
			break
		}

		gitFiles = append(gitFiles, file.Name)
	}

	return gitFiles, nil
}

func listOsFiles(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})

	return files, err
}
