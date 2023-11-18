package main

import "github.com/go-git/go-git/v5"

func openRepo(gitdir *string) (repo *git.Repository, tree *git.Worktree, err error) {
	repo, err = git.PlainOpen(*gitdir)
	if err != nil {
		return
	}
	tree, err = repo.Worktree()
	return
}
