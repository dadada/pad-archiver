package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"
	"crypto/tls"
	"sync"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

var commitmu sync.Mutex

func commit(
	tree *git.Worktree,
	padfile string,
	url string,
) (plumbing.Hash, error) {
	commitmu.Lock()
	defer commitmu.Unlock()

	commit, err := tree.Commit(
		fmt.Sprintf("Updated %s from %s", padfile, url),
		&git.CommitOptions{
			Author: &object.Signature {
				Name: "Pad Archiver",
				Email: "pad-archiver@dadada.li",
				When: time.Now(),
			},
		},
	)

	if err != nil {
		return plumbing.ZeroHash, fmt.Errorf("Failed to commit %s: %w", padfile, err)
	}

	return commit, nil
}

func update(
	tree *git.Worktree,
	url string,
) (plumbing.Hash, error) {
	res, err := http.Get(url + "/export/txt")
	if err != nil {
		return plumbing.ZeroHash, fmt.Errorf("Failed to get pad at %s: %w", url, err)
	}

	defer res.Body.Close()

	padfile := path.Base(url) + ".txt"

	padpath := filepath.Join(tree.Filesystem.Root(), padfile)
	out, err := os.Create(padpath)

	written, err := io.Copy(out, res.Body)
	if err != nil {
		return plumbing.ZeroHash, fmt.Errorf("Failed to write pad to file at %s: %w", padfile, err)
	}

	if written < 100 {
		return plumbing.ZeroHash, fmt.Errorf("Skipping update of %s, because pad has likely been removed from %s", padfile, url)
	}

	status, err := tree.Status()
	if err != nil {
		return plumbing.ZeroHash, fmt.Errorf("Failed to get status of %s: %w", padfile, err)
	}

	if status.IsClean() {
		return plumbing.ZeroHash, fmt.Errorf("No changes recorded for %s", url)
	}

	if _, err = tree.Add(padfile); err != nil {
		return plumbing.ZeroHash, fmt.Errorf("Failed to stage %s: %w", padfile, err)
	}

	return commit(tree, padfile, url)
}

func main() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get working directory %s", err)
	}

	gitdir := flag.String(
		"C",
		cwd,
		"git directory",
	)
	flag.Parse()

	repo, err := git.PlainOpen(*gitdir)
	if err != nil {
		log.Fatalf("Failed to open git repo %s: %s", *gitdir, err)
	}

	tree, err := repo.Worktree()
	if err != nil {
		log.Fatalf("Failed to open git worktree %s", err)
	}

	scanner := bufio.NewScanner(os.Stdin)

	var wg sync.WaitGroup
	for scanner.Scan() {
		wg.Add(1)
		padurl := scanner.Text()

		go func() {
			defer wg.Done()
			if _, err := update(tree, padurl); err != nil {
				log.Printf("%s", err)
			} else {
				log.Printf("Updated %s", padurl)
			}
		}()
	}
	wg.Wait()

	tree.Clean(&git.CleanOptions{})
}
