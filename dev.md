# Development
Developers log for todos and plans.


## Todo
- Move error handling from execute.go to calls in cmd/?
- Change updating print so that if file is being added for the first time, it says "adding" rather than "updating".
- Check if includes/excludes works fine if is directory.
- Remove trailing slash from SystemFileName, which occurs whenever `path` is passed with no trailing slash.
- Plan out project structure in terms of packages, package inputs/outputs, etc.
- Add warning if an excluded file is still stored in target repository.
- Try and improve handling of file walking limit.
  - Replace walk() with walkDir() for efficiency?
- Add warning before syncing; list all files that are about to be synced, then wait for user confirmation.

### Resources
- https://github.com/go-git/go-git
- https://github.com/sergi/go-diff
- https://github.com/go-yaml/yaml


## Internal
### Config File
- [x] Searched for in multiple file locations (~/.config/docon/ and current directory).
- [x] Flag for specifying config file location.
- [ ] Config field for renaming files; dictionary where keys are original file names and values are new file names. Files are renamed after being synced and before being committed.

### Interactive Mode
- [ ] Adds extra step where all files are first listed with a numbered index.
- [ ] User provides index number(s) for the files they wish to sync.


## Commands
### Sync
- [x] Sync system files with target files.
- [x] Syncs all files in config by default
- [x] Flag for automatically running commit command after syncing.
- [ ] Flag for specifying file permissions for when new folders/files are created automatically.
- [ ] Flag for verbose mode which prints out last modified, size, etc of each file being processed.
- [ ] Flag for specifying which dotfiles to sync.
- [ ] Flag for specifying which dotfiles to ignore.
- [ ] Flag for interactive mode.

### Diff
- [x] Show diff between system dotfiles and repo dotfiles.
- [x] Diff is shown for all files in config by default.
- [x] Diff is coloured and similar to Git diff.
- [ ] Flag for specifying which dotfiles to show diff for.
- [ ] Flag for specifying which dotfiles to ignore.
- [ ] Flag for interactive mode.

### Commit
- [x] Automatically commits files/directories to Git.
- [x] Commits all files in config by default.
- [x] Uses a default commit message template if none is provided.
- [x] Config fields accept keywords.
- [x] Config field for specifying commit messages for each config group.
- [x] Config field for specifying global commit message for every config group (overriden for groups where individual messages are provided).
- [x] Config field for author and email. If not provided, uses git config file by default.
- [ ] Config field for adding either entire directory in a single commit (only if >1 file), or individual files.
- [ ] Flag for verbose mode which prints out commit object.
- [ ] Flag for specifying which dotfiles to commit.
- [ ] Flag for specifying which dotfiles to ignore.
- [ ] Flag for interactive mode.

### Check
- [ ] Checks if dotfiles need to be synced or if they're up to date.
- [ ] Check all files in config by default.
- [ ] Show which files need updating.
- [ ] Notifies last time system was synced.
- [ ] Notifies if target dotfiles don't exist/cannot be found.
- [ ] Files list includes the last modified timestamp for both system dotfiles and target dotfiles.
- [ ] Files list includes line count difference.
- [ ] Flag for specifying which dotfiles to commit.
- [ ] Flag for specifying which dotfiles to ignore.
- [ ] Flag for interactive mode.

### Backup
- [ ] Backup dotfiles to specified file.
- [ ] Backup all files in config by default.
- [ ] User specifies location and extension.
- [ ] Flag for specifying which dotfiles to backup.
- [ ] Flag for specifying which dotfiles to ignore.
- [ ] Flag for interactive mode.

### Rollback
- [ ] Rollback system dotfiles to same state as target dotfiles.
- [ ] Rollback all files in config by default.
- [ ] Notifies which files will be rolled back before running and asks for confirmation.
- [ ] A folder containing all the original system dotfiles is created as a backup by default.
- [ ] Flag for disabling default backup.
- [ ] Flag for specifying backup folder location and extension.
- [ ] Flag for specifying which dotfiles to backup.
- [ ] Flag for specifying which dotfiles to ignore.
- [ ] Flag for interactive mode.

### Restore
- [ ] Uses a backup folder to restore a system to that state.

### Verify
- [x] Check config file syntax and paths.
- [ ] Notify if excluded files are already present in repo.

### Pkg
- [x] Generates a packages list in target directory.
- [x] Package list is generated based on Linux distribution being used.
- [ ] Includes "Added", "Removed", and "Updated" sections when printing differences.
- [ ] Config field for specifying package list name and location.
- [ ] Flag for removing version number from pkglist entries.
- [ ] Flag for quiet mode where no output is printed.
- [ ] Flag for including package history which, instead of replacing previous file, creates a new one with the date in file name (e.g. pkglist-26.06.21).
