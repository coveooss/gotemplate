package template

import (
	"fmt"

	git "gopkg.in/src-d/go-git.v4"
)

const (
	gitBase = "Git functons"
)

var gitFuncs = dictionary{
	"current_branch": currentBranch,
	"current_commit": currentCommit,
	"origin":         origin,
}

var gitFuncsArgs = arguments{
	"current_branch": {"path"},
	"current_commit": {"path"},
	"origin":         {"path"},
}

var gitFuncsAliases = aliases{}

var gitFuncsHelp = descriptions{}

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
		return "", fmt.Errorf("Not currently in a branch: %s", head.String())
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
