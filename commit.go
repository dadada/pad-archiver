package main

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

var (
	cm          sync.Mutex
	nothingToDo = errors.New("Nothing to do for unmodified file")
)

func commit(
	tree *git.Worktree,
	padfile string,
	url string,
) (commit plumbing.Hash, err error) {
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

	commit, err = tree.Commit(
		fmt.Sprintf("Updated %s from %s", padfile, url),
		&git.CommitOptions{
			All: false,
			Author: &object.Signature{
				Name:  "Pad Archiver[bot]",
				Email: "noreply@example.com",
				When:  time.Now(),
			},
		},
	)

	if err != nil {
		return plumbing.ZeroHash, fmt.Errorf("Failed to commit %s: %w", padfile, err)
	}

	return
}
