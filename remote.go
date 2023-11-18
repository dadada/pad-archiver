package main

import (
	"log"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
)

const (
	defaultRemoteName = "pad-archiver"
)

func auth(username *string, password *string) (auth *githttp.BasicAuth) {
	return &githttp.BasicAuth{
		Username: *username,
		Password: *password,
	}
}

func pushRepo(repo *git.Repository, remoteUrl *string, auth *githttp.BasicAuth) (err error) {
	if _, err = createRemote(repo, defaultRemoteName, *remoteUrl); err != nil {
		return
	}
	if err = push(auth, repo, defaultRemoteName); err != nil {
		if err == git.NoErrAlreadyUpToDate {
			log.Println("Already up-to-date")
		} else {
			return
		}
	} else {
		log.Println("Pushed changes to remote")
	}
	return
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
