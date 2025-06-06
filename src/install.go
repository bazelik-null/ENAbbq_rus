package main

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
)

func install(path string) error {
	appDir, _ := os.Getwd()
	metaPath := filepath.Join(appDir, "resources", "meta.json")

	// Open meta.json
	data, err := os.ReadFile(metaPath)
	if err != nil {
		return err
	}

	// Extract data
	var meta struct {
		Files []struct {
			Src string `json:"src"`
			Dst string `json:"dst"`
		} `json:"files"`
	}
	json.Unmarshal(data, &meta)

	// Copy all files
	for _, file := range meta.Files {
		src := filepath.Join(appDir, file.Src)
		dst := filepath.Join(path, file.Dst)

		err := copyFile(src, dst)
		if err != nil {
			return err
		}
	}

	return nil
}

func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	return nil
}
