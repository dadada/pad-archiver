package main

import (
	"bufio"
	"crypto/tls"
	"errors"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/go-git/go-git/v5"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
)

const (
	defaultRemoteName = "pad-archiver"
)

var (
	nothingToDo = errors.New("Nothing to do for unmodified file")
)

func main() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	gitdir, doPush, username, password, remoteUrl := getArgs()

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
		pushRepo(repo, remoteUrl, auth)
	}
}
