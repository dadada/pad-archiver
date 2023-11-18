package main

import (
	"flag"
	"os"
)

func getArgs(workingdir *string) (gitdir *string, doPush *bool, username *string, password *string, remoteUrl *string) {
	gitdir = flag.String(
		"C",
		*workingdir,
		"The directory containing the git repository in which to archive the pads.",
	)
	doPush = flag.Bool(
		"push",
		false,
		"Push the changes to the remote specified by remoteUrl.",
	)
	username = flag.String(
		"username",
		"",
		"The username for authenticating to the remote.",
	)
	password = flag.String(
		"password",
		os.Getenv("GIT_PASSWORD"),
		"The password for authenticating to the remote. Can also be specified via the environment variable GIT_PASSWORD.",
	)
	remoteUrl = flag.String(
		"url",
		"",
		"URL to push changes to.",
	)

	flag.Parse()

	return
}
