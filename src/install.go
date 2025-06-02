package main

import (
	"encoding/json"
	//"fmt"
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

	// WORK IN PROGRESS
	/*	for i := 0; i <= 11; i++ {
		level := "level" + fmt.Sprint(i)
		srcLevel := filepath.Join(appDir, "resources", "levels", level)
		dstLevel := filepath.Join(path, "ENA-4-DreamBBQ_Data", level)

		err := copyFile(srcLevel, dstLevel)
		if err != nil {
			return err
		}
	} */
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
