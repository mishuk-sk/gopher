package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
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
			path: "tmp",
			files: []string{"birthday_001.txt",
				"birthday_002.txt",
				"birthday_003.txt",
				"christmas 2016 (1 of 100).txt",
				"christmas 2016 (2 of 100).txt",
				"christmas 2016 (3 of 100).txt"},
		},
	}
	tmpDir, err := prepareTestDirTree(directories)
	if err != nil {
		t.Fatalf("Can't create temp directory tree: %v\n", err)
	}
	defer os.RemoveAll(tmpDir)
	os.Chdir(tmpDir)
	if err := filepath.Walk(".", walkFunc(rename)); err != nil {
		t.Fatalf("Error during walking. Err: %v\n", err)
	}
	if err := filepath.Walk(".", outputWalk(t)); err != nil {
		t.Fatalf("Error during walking. Err: %v\n", err)
	}
}

func outputWalk(t *testing.T) filepath.WalkFunc {
	var reg = regexp.MustCompile(`(.+?)_([0-9]+)(.+)`)
	return func(path string, info os.FileInfo, err error) error {
		if reg.MatchString(info.Name()) {
			t.Fatalf("Unexpected match of filename %s. Expected to be raplaced.", info.Name())
		}
		return nil
	}
}
