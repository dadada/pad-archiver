package main

import (
	"bufio"
	"crypto/tls"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sync"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
)

const (
	defaultRemoteName = "pad-archiver"
)

var (
	nothingToDo = errors.New("Nothing to do for unmodified file")
)

var cm sync.Mutex

func commit(
	tree *git.Worktree,
	padfile string,
	url string,
) (plumbing.Hash, error) {
	cm.Lock()
	defer cm.Unlock()

	if _, err := tree.Add(padfile); err != nil {
		return plumbing.ZeroHash, fmt.Errorf("Failed to stage %s: %w", padfile, err)
	}

	status, err := tree.Status()
	if err != nil {
		return plumbing.ZeroHash, fmt.Errorf("Failed to get status of %s", padfile)
	}

	fileStatus := status.File(padfile)
	if fileStatus.Staging != git.Added && fileStatus.Staging != git.Modified {
		return plumbing.ZeroHash, nothingToDo
	}

	commit, err := tree.Commit(
		fmt.Sprintf("Updated %s from %s", padfile, url),
		&git.CommitOptions{
			All: false,
			Author: &object.Signature{
				Name:  "Pad Archiver",
				Email: "pad-archiver@dadada.li",
				When:  time.Now(),
			},
		},
	)

	if err != nil {
		return plumbing.ZeroHash, fmt.Errorf("Failed to commit %s: %w", padfile, err)
	}

	return commit, nil
}

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

func createRemote(
	repo *git.Repository,
	remote string,
	url string,
) (*git.Remote, error) {
	newRemote, err := repo.Remote(remote)
	if err != nil {
		log.Printf("Creating new git remote %s with URL %s", remote, url)
		return repo.CreateRemote(&config.RemoteConfig{
			Name: remote,
			URLs: []string{url},
		})
	} else {
		log.Printf("Using remote %s with URL %s", remote, url)
	}

	return newRemote, nil
}

func push(
	auth *githttp.BasicAuth,
	r *git.Repository,
	remote string,
) error {
	return r.Push(&git.PushOptions{
		RemoteName: remote,
		Auth:       auth,
	})
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
		"The directory containing the git repository in which to archive the pads.",
	)
	doPush := flag.Bool(
		"push",
		false,
		"Push the changes to the remote specified by remoteUrl.",
	)
	username := flag.String(
		"username",
		"",
		"The username for authenticating to the remote.",
	)
	password := flag.String(
		"password",
		os.Getenv("GIT_PASSWORD"),
		"The password for authenticating to the remote. Can also be specified via the environment variable GIT_PASSWORD.",
	)
	remoteUrl := flag.String(
		"url",
		"",
		"URL to push changes to.",
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

	defer tree.Clean(&git.CleanOptions{Dir: true})

	filesystemRoot := tree.Filesystem.Root()
	scanner := bufio.NewScanner(os.Stdin)

	var wg sync.WaitGroup
	for scanner.Scan() {
		wg.Add(1)
		padurl := scanner.Text()

		go func() {
			defer wg.Done()
			padfile, err := download(filesystemRoot, padurl)
			if err != nil {
				log.Printf("%s", err)

				return
			}
			log.Printf("Downloaded %s", padurl)
			if _, err := commit(tree, padfile, padurl); err != nil {
				if err == nothingToDo {
					log.Printf("Nothing to do for %s", padfile)
				} else {
					log.Printf("%s", err)
				}
			} else {
				log.Printf("Committed %s", padfile)
			}
		}()
	}
	wg.Wait()

	auth := &githttp.BasicAuth{
		Username: *username,
		Password: *password,
	}

	if *doPush == true {
		if _, err := createRemote(repo, defaultRemoteName, *remoteUrl); err != nil {
			log.Fatalf("%s", err)
		}
		if err := push(auth, repo, defaultRemoteName); err != nil {
			if err == git.NoErrAlreadyUpToDate {
				log.Println("Already up-to-date")
			} else {
				log.Fatalf("%s", err)
			}
		} else {
			log.Println("Pushed changes to remote")
		}
	}
}
