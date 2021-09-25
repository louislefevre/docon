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

const (
	gitStatusKeyword keywordType = "{status}"
	gitFileKeyword   keywordType = "{file}"
	gitDirKeyword    keywordType = "{dir}"
	gitUserKeyword   keywordType = "{user}"
	gitEmailKeyword  keywordType = "{email}"
)

var gitKeywords = keywordSet{
	newKeyword(gitStatusKeyword),
	newKeyword(gitFileKeyword),
	newKeyword(gitDirKeyword),
	newKeyword(gitUserKeyword),
	newKeyword(gitEmailKeyword),
}

func CommitAll(config *configuration) error {
	for _, group := range config.Sources {
		for _, df := range group.dotfiles {
			var err error
			if group.CommitMsg != "" {
				err = commit(config.TargetPath, df.targetFile.name, group.CommitMsg)
			} else if config.Git.CommitMsg != "" {
				err = commit(config.TargetPath, df.targetFile.name, config.Git.CommitMsg)
			} else {
				err = commit(config.TargetPath, df.targetFile.name, "{Status} {file}")
			}
			if err != nil {
				return newError(err, fmt.Sprintf("Failed to commit %s", df.targetFile.name))
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

	status, err := tree.Status()
	if err != nil {
		return err
	}

	fileStatus := status.File(file).Staging
	if fileStatus == git.Untracked {
		return nil
	}

	statusString, err := getStatusString(fileStatus)
	if err != nil {
		return err
	}

	author, err := getAuthorSignature()
	if err != nil {
		return err
	}

	for _, kw := range gitKeywords {
		switch kw.kwType {
		case gitStatusKeyword:
			msg = kw.transform(msg, statusString)
		case gitFileKeyword:
			msg = kw.transform(msg, file)
		case gitDirKeyword:
			msg = kw.transform(msg, "")
		case gitUserKeyword:
			msg = kw.transform(msg, author.Name)
		case gitEmailKeyword:
			msg = kw.transform(msg, author.Email)
		}
	}

	_, err = tree.Commit(msg, &git.CommitOptions{Author: author})
	if err != nil {
		return err
	}

	return nil
}

func getStatusString(status git.StatusCode) (string, error) {
	switch status {
	case git.Added:
		return "add", nil
	case git.Modified:
		return "update", nil
	case git.Deleted:
		return "delete", nil
	case git.Renamed:
		return "rename", nil
	case git.Copied:
		return "copy", nil
	default:
		return "", fmt.Errorf("unknown git status code %s", string(status))
	}
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
