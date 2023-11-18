package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

func download(
	gitdir string,
	url string,
) (string, error) {
	res, err := http.Get(url + "/export/txt")
	if err != nil {
		return "", fmt.Errorf("Failed to get pad at %s: %w", url, err)
	}

	defer res.Body.Close()

	padfile := path.Base(url) + ".txt"

	padpath := filepath.Join(gitdir, padfile)
	out, err := os.Create(padpath)

	written, err := io.Copy(out, res.Body)
	if err != nil {
		return "", fmt.Errorf("Failed to write pad to file at %s: %w", padfile, err)
	}

	if written < 100 {
		return "", fmt.Errorf("Skipping update of %s, because pad has likely been removed from %s", padfile, url)
	}

	return padfile, nil
}
