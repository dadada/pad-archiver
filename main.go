package main

import (
	"bufio"
	"crypto/tls"
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
	"errors"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
)

const DefaultRemoteName = "pad-archiver"

var (
	NothingToDo = errors.New("Nothing to do for unmodified file")
)

var Commitmu sync.Mutex


func Commit(
	tree *git.Worktree,
	padfile string,
	url string,
) (plumbing.Hash, error) {
	Commitmu.Lock()
	defer Commitmu.Unlock()

	if _, err := tree.Add(padfile); err != nil {
		return plumbing.ZeroHash, fmt.Errorf("Failed to stage %s: %w", padfile, err)
	}

	status, err := tree.Status()
	if err != nil {
		return plumbing.ZeroHash, fmt.Errorf("Failed to get status of %s", padfile)
	}

	fileStatus := status.File(padfile)
	if fileStatus.Staging != git.Added && fileStatus.Staging != git.Modified {
		return plumbing.ZeroHash, NothingToDo
	}

	commit, err := tree.Commit(
		fmt.Sprintf("Updated %s from %s", padfile, url),
		&git.CommitOptions{
			All: false,
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

func Download(
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


func CreateRemote(
	repo *git.Repository,
	remote string,
) (*git.Remote, error) {
	newRemote, err := repo.Remote(DefaultRemoteName)
	if err != nil {
		log.Printf("Creating new git remote %s with URL %s", DefaultRemoteName, remote)
		return repo.CreateRemote(&config.RemoteConfig{
			Name: DefaultRemoteName,
			URLs: []string{remote},
		})
	} else {
		log.Printf("Using remote %s with URL %s", DefaultRemoteName, remote)
	}

	return newRemote, nil
}

func Push(
	auth *githttp.BasicAuth,
	remote *git.Remote,
) error {
	return remote.Push(&git.PushOptions{
		Auth: auth,
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
		"git directory",
	)
	push := flag.Bool(
		"push",
		false,
		"push repository to remote",
	)
	username := flag.String(
		"username",
		"",
		"username",
	)
	password := flag.String(
		"password",
		"",
		"password",
	)
	remoteUrl := flag.String(
		"remote",
		"",
		"remote",
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

	filesystemRoot := tree.Filesystem.Root()
	scanner := bufio.NewScanner(os.Stdin)

	var wg sync.WaitGroup
	for scanner.Scan() {
		wg.Add(1)
		padurl := scanner.Text()

		go func() {
			defer wg.Done()
			padfile, err := Download(filesystemRoot, padurl)
			if err != nil {
				log.Printf("%s", err)

				return
			}
			log.Printf("Downloaded %s", padurl)
			if _, err := Commit(tree, padfile, *remoteUrl); err != nil {
				if err == NothingToDo {
					log.Printf("Nothing to do for %s", padfile)
				} else {
					log.Fatalf("%s", err)
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

	remote, err := CreateRemote(repo, *remoteUrl)
	if err != nil {
		log.Fatalf("%s", err)
	}

	if *push == true {
		if err := Push(auth, remote); err != nil {
			if err == git.NoErrAlreadyUpToDate {
				log.Println("Already up-to-date")
			} else {
				log.Fatalf("%s", err)
			}
		} else {
			log.Println("Pushed changes to remote")
		}
	}

	tree.Clean(&git.CleanOptions{Dir: true})
}
