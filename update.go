package main

import (
	"bufio"
	"log"
	"sync"

	"github.com/go-git/go-git/v5"
)

func updatePads(pads *bufio.Scanner, tree *git.Worktree) {
	defer tree.Clean(&git.CleanOptions{Dir: true})

	filesystemRoot := tree.Filesystem.Root()

	var wg sync.WaitGroup
	for pads.Scan() {
		wg.Add(1)
		padurl := pads.Text()

		go func() {
			defer wg.Done()
			if _, err := updatePad(filesystemRoot, padurl, tree); err != nil {
				log.Printf("%s", err)
			}
		}()
	}
	wg.Wait()
}

func updatePad(filesystemRoot string, padurl string, tree *git.Worktree) (padfile string, err error) {
	padfile, err = download(filesystemRoot, padurl)
	if err != nil {
		return
	}
	log.Printf("Downloaded %s", padurl)
	if _, err = commit(tree, padfile, padurl); err != nil {
		return
	} else {
		log.Printf("Committed %s", padfile)
		return
	}
}
