package main

import (
	"fmt"
	"os"
)

func main() {

}

func walk(path string, info os.FileInfo, err error) error {
	fmt.Printf("visited file or directory: %s\n", path)
	if err != nil {
		fmt.Printf("Visited with error: %v\n", err)
	}
	if !info.IsDir() {
		if err := rename(path, "file.txt"); err != nil {
			panic(err)
		}
	}
	return nil
}

func rename(old, new string) error {
	return os.Rename(old, new)
}
