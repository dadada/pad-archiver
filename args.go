package main

import (
	"flag"
	"log"
	"os"
)

func getArgs() (*string, *bool, *string, *string, *string) {
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

	return gitdir, doPush, username, password, remoteUrl
}
