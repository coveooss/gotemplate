package template

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

const (
	branchName  = "test-branch"
	dummyRemote = "git@github.com/dummy/remote"
)

func TestGetCurrentBranch(t *testing.T) {
	path, _, _ := initTestRepository(t)
	defer os.RemoveAll(path)

	result, err := currentBranch(path)

	assert.NoError(t, err)
	assert.Equal(t, branchName, result)
}

func TestGetCurrentCommit(t *testing.T) {
	path, _, commit := initTestRepository(t)
	defer os.RemoveAll(path)

	result, err := currentCommit(path)

	assert.NoError(t, err)
	assert.Equal(t, commit.Hash.String(), result)
}

func TestOrigin(t *testing.T) {
	path, _, _ := initTestRepository(t)
	defer os.RemoveAll(path)

	result, err := origin(path)

	assert.NoError(t, err)
	assert.Equal(t, dummyRemote, result)
}

func initTestRepository(t *testing.T) (string, *git.Repository, *object.Commit) {
	path := t.TempDir()

	// Init the repository
	repo, err := git.PlainInit(path, false)
	assert.NoError(t, err)
	worktree, err := repo.Worktree()
	assert.NoError(t, err)

	// Create a file to commit
	filename := filepath.Join(path, "example-git-file")
	err = os.WriteFile(filename, []byte("hello world!"), 0644)
	assert.NoError(t, err)
	_, err = worktree.Add("example-git-file")
	assert.NoError(t, err)

	// Create a commit
	commit, err := worktree.Commit("example go-git commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "John Doe",
			Email: "john@doe.org",
			When:  time.Now(),
		},
	})
	assert.NoError(t, err)
	commitObject, err := repo.CommitObject(commit)
	assert.NoError(t, err)

	// Create a branch
	headRef, err := repo.Head()
	assert.NoError(t, err)
	ref := plumbing.NewHashReference("refs/heads/"+branchName, headRef.Hash())
	err = repo.Storer.SetReference(ref)
	assert.NoError(t, err)
	err = worktree.Checkout(&git.CheckoutOptions{Create: false, Force: false, Branch: ref.Name()})
	assert.NoError(t, err)

	// Set a remote
	repo.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{dummyRemote},
	})

	return path, repo, commitObject
}
