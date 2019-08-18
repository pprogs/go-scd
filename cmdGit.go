package main

import (
	"fmt"
	"os"

	"gopkg.in/src-d/go-git.v4/plumbing"

	git "gopkg.in/src-d/go-git.v4"
	git_http "gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

func gitCommand(c *command) {

	if ok, _ := PathExists(c.LocalRepo); ok {
		pull(c)
	} else {
		clone(c)
	}
}

func clone(c *command) {

	branch := fmt.Sprintf("refs/heads/%s", c.Branch)
	b := plumbing.ReferenceName(branch)

	opts := &git.CloneOptions{
		Auth: &git_http.BasicAuth{
			Username: "login",
			Password: c.Token,
		},
		URL:           c.RemoteRepo,
		Progress:      os.Stdout,
		ReferenceName: b,
	}

	_, err := git.PlainClone(c.LocalRepo, false, opts)

	lg.Printf("%v\n", err)
}

func pull(c *command) {

	// We instance\iate a new repository targeting the given path (the .git folder)
	r, err := git.PlainOpen(c.LocalRepo)
	CheckIfError(err)

	// Get the working directory for the repository
	w, err := r.Worktree()
	CheckIfError(err)

	opts := &git.PullOptions{
		Auth: &git_http.BasicAuth{
			Username: "login", // yes, this can be anything except an empty string
			Password: c.Token,
		},
		RemoteName: "origin",
		Progress:   os.Stdout,
	}

	// Pull the latest changes from the origin remote and merge into the current branch
	err = w.Pull(opts)
	if err == git.NoErrAlreadyUpToDate {
		lg.Printf("Already up-to-date\n")
		return
	}

	CheckIfError(err)
}
