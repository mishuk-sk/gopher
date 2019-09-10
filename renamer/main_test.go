package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

type directory struct {
	path  string
	files []string
}

func prepareTestDirTree(dirs []directory) (string, error) {
	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		return "", fmt.Errorf("error creating temp dir: %v", err)
	}
	for _, dir := range dirs {
		dirpath := filepath.Join(tmpDir, dir.path)
		err = os.MkdirAll(dirpath, 0755)
		if err != nil {
			return "", fmt.Errorf("error creating directory (%s). Err:%v", dir.path, err)
		}
		for _, filename := range dir.files {
			f, err := os.Create(filepath.Join(dirpath, filename))
			if err != nil {
				return "", fmt.Errorf("error creating file %s. Err:%v", filename, err)
			}
			f.Close()
		}
	}
	return tmpDir, nil
}
func TestWalking(t *testing.T) {
	directories := []directory{
		directory{
			path:  "tmp",
			files: []string{"top.txt", "bot.txt"},
		},
	}
	tmpDir, err := prepareTestDirTree(directories)
	if err != nil {
		t.Fatalf("Can't create temp directory tree: %v\n", err)
	}
	defer os.RemoveAll(tmpDir)
	os.Chdir(tmpDir)
	if err := filepath.Walk(".", walk); err != nil {
		t.Fatalf("Error during walking. Err: %v\n", err)
	}
	if err := filepath.Walk(".", outputWalk); err != nil {
		t.Fatalf("Error during walking. Err: %v\n", err)
	}
}

func outputWalk(path string, info os.FileInfo, err error) error {
	fmt.Printf("visited file or directory: %s\n", path)
	return nil
}
