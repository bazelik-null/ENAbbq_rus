package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func getURL() (string, int) {
	client := &http.Client{}

	repoOwner := "bazelik-null"
	repoName := "ENAbbq_rus"

	req, _ := http.NewRequest("GET", fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", repoOwner, repoName), nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3")

	resp, err := client.Do(req)
	if err != nil {
		return "", http.StatusInternalServerError
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", resp.StatusCode
	}

	var release struct {
		Assets []struct {
			BrowserDownloadURL string `json:"browser_download_url"`
			Name               string `json:"name"`
		} `json:"assets"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&release)

	// Search for assets archive
	for _, asset := range release.Assets {
		if asset.Name == "resources.zip" {
			return asset.BrowserDownloadURL, http.StatusOK
		}
	}

	return "", http.StatusNotFound
}

func downloadFile(url, filename string) int {
	resp, err := http.Get(url)
	if err != nil {
		return http.StatusInternalServerError
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return resp.StatusCode
	}

	f, err := os.Create(filename)
	if err != nil {
		return http.StatusInternalServerError
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return http.StatusInternalServerError
	}

	return 0
}

func download() int {
	url, code := getURL()
	if code != http.StatusOK {
		return code
	}

	// Check if URL is available
	if _, err := http.Head(url); err != nil {
		return http.StatusBadGateway
	}

	filename := "resources.zip"
	code = downloadFile(url, filename)
	if code != 0 {
		return code
	}

	r, _ := zip.OpenReader(filename)
	defer r.Close()

	// Extract zip archive
	for _, f := range r.File {
		rc, _ := f.Open()
		defer rc.Close()

		fpath := filepath.Join(".", f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, f.Mode())
		} else {
			f, _ := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			defer f.Close()

			io.Copy(f, rc)
		}
	}

	// Delete temporary zip archive
	_ = os.Remove(filename)
	return 0
}
