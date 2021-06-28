package internal

import (
	"fmt"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/tcnksm/go-gitconfig"
)

// TODO: Add config file field for adding either entire directory (only if >1 file) or individual files.
// TODO: Add author section to config file and check if empty before running this.
// TODO: Add check for empty commit message: if empty, use default. Maybe "{empty}" keyword in config if
//       user wants to keep it truly empty?
// TODO: Add check to see whether file is being added, updated, or removed.

func commitAll(config *configuration) error {
	for _, group := range config.Sources {
		for _, df := range group.dotfiles {
			var err error
			if group.CommitMsg != "" {
				err = commit(config.TargetPath, df.targetFile.name, group.CommitMsg)
			} else if config.Git.CommitMsg != "" {
				err = commit(config.TargetPath, df.targetFile.name, config.Git.CommitMsg)
			} else {
				msg := fmt.Sprintf("Update %s", df.targetFile.name)
				err = commit(config.TargetPath, df.targetFile.name, msg)
			}
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func commit(dir string, file string, msg string) error {
	repo, err := git.PlainOpen(dir)
	if err != nil {
		return err
	}

	tree, err := repo.Worktree()
	if err != nil {
		return err
	}

	_, err = tree.Add(file)
	if err != nil {
		return err
	}

	author, err := getAuthorSignature()
	if err != nil {
		return err
	}

	_, err = tree.Commit(msg, &git.CommitOptions{Author: author})
	if err != nil {
		return err
	}

	return nil
}

func getAuthorSignature() (*object.Signature, error) {
	author := &object.Signature{}

	if username, err := gitconfig.Username(); err == nil {
		author.Name = username
	} else {
		return author, err
	}

	if email, err := gitconfig.Email(); err == nil {
		author.Email = email
	} else {
		return author, err
	}

	author.When = time.Now()

	return author, nil
}
