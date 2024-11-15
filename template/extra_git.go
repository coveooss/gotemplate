package template

import (
	"fmt"

	git "github.com/go-git/go-git/v5"
)

const (
	gitBase = "Git functions"
)

var gitFuncs = dictionary{
	"currentBranch": currentBranch,
	"currentCommit": currentCommit,
	"origin":        origin,
}

var gitFuncsArgs = arguments{
	"currentBranch": {"path"},
	"currentCommit": {"path"},
	"origin":        {"path"},
}

var gitFuncsAliases = aliases{}

var gitFuncsHelp = descriptions{
	"currentBranch": "Returns the name of the currently checked out git branch at the given path",
	"currentCommit": "Returns the hash of the currently checked out git commit at the given path",
	"origin":        "Returns the git origin remote URL at the given path",
}

func (t *Template) addGitFuncs() {
	t.AddFunctions(gitFuncs, gitBase, FuncOptions{
		FuncHelp:    gitFuncsHelp,
		FuncArgs:    gitFuncsArgs,
		FuncAliases: gitFuncsAliases,
	})
}

func currentCommit(path string) (string, error) {
	repository, err := git.PlainOpen(path)
	if err != nil {
		return "", err
	}
	head, err := repository.Head()
	if err != nil {
		return "", err
	}
	return head.Hash().String(), nil
}

func currentBranch(path string) (string, error) {
	repository, err := git.PlainOpen(path)
	if err != nil {
		return "", err
	}
	head, err := repository.Head()
	if err != nil {
		return "", err
	}
	if !head.Name().IsBranch() {
		return "", fmt.Errorf("not currently in a branch: %s", head.String())
	}
	return head.Name().Short(), nil
}

func origin(path string) (string, error) {
	repository, err := git.PlainOpen(path)
	if err != nil {
		return "", err
	}
	remote, err := repository.Remote("origin")
	if err != nil {
		return "", err
	}
	return remote.Config().URLs[0], nil
}
