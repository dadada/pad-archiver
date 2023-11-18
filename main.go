package main

import (
	"bufio"
	"crypto/tls"
	"log"
	"net/http"
	"os"
)

func main() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	workingdir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get working directory %s", err)
	}

	gitdir, doPush, username, password, remoteUrl := getArgs(&workingdir)

	repo, tree, err := openRepo(gitdir)
	if err != nil {
		if repo == nil {
			log.Fatalf("Failed to open git repo %s: %s", *gitdir, err)
		}
		if tree == nil {
			log.Fatalf("Failed to open git worktree %s", err)
		}
	}

	padstxt := bufio.NewScanner(os.Stdin)

	updatePads(padstxt, tree)

	auth := auth(username, password)

	if *doPush == true {
		if err := pushRepo(repo, remoteUrl, auth); err != nil {
			log.Fatalf("Failed to push repo: %s", err)
		}
	}
}
