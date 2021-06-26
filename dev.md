# Development
Developers log for todos and plans.


## Todo
- Create a Path struct.
- Create a filesystem package with tools for manipulating file system.
- Rename cmd files to end with "Cmd".
- Move error handling from execute.go to calls in cmd/?
- Split up repo list into a separate file? Have config exclusively for docon configuration options?
- Change updating print so that if file is being added for the first time, it says "adding" rather than "updating".
- Add check for if passed repo is directory and exists.
- Check if includes/excludes works fine if is directory.
- Remove trailing slash from SystemFileName, which occurs whenever `path` is passed with no trailing slash.
- Plan out project structure in terms of packages, package inputs/outputs, etc.
- Add custom error struct with pretty print method.
  - Struct includes error message from the system, simpler error message written by me, and perhaps code failure line/file if possible.

### Resources
- https://github.com/go-git/go-git
- https://github.com/sergi/go-diff
- https://github.com/go-yaml/yaml


## Internal
### Config File
- File locations (priority order):
  - User passed directory (via flag)
  - Current working directory (looks for docon/config.yaml)
  - .config/docon/config.yaml
  - /etc/docon/config.yaml


## Commands
### Sync
- Syncs system files with target files.
- Default syncs all files in config.
- Flag for specifying which dotfiles to sync.
- Flag for specifying which dotfiles to ignore (all others will be synced).
- Flag for interactive sync:
  - Adds extra step where all files are first listed with a numbered index.
  - User provides index number(s) for the files they wish to sync.
- Can automatically run commit command after syncing with flag.
- Flag for specifying file permissions for when new folders/files are created automatically.
- Flag for verbose mode which prints out last modified, size, etc of each file being processed.


### Diff
- Show diff between system dotfiles and repo dotfiles.
- Diff is coloured and similar to Git diff.
- Default shows diff between all dotfiles.
- Flag for specifying which dotfiles to show diff for.
- Flag for interactive diff:
  - Adds extra step where all files are first listed with a numbered index.
  - User provides index number(s) for the files they wish to show the diff.

### Commit
- Automatically commits files/directories to Git.
- Automatically generates git commit messages:
  - Default is "Update directory/filename".
  - Can specify specific messages for each config group in config.yaml file.

### Check
- Checks if dotfiles need to be synced or if they're up to date.
- If sync is needed, show which files need updating. If not, do nothing.
- Files list includes the last modified timestamp for both system dotfiles and target dotfiles (if exists).
- Files list also includes line count difference.
- Outputs last time system was synced using time-stamp.

### Backup
- Backup dotfiles to location.
- User specifies location and extension.
- Default will backup all dotfiles.
- Flag for specifying which dotfiles to backup.
- Flag for specifying which dotfiles to ignore (all others will be backed up).
- Flag for interactive backup:
  - Adds extra step where all files are first listed with a numbered index.
  - User provides index number(s) for the files they wish to backup.

### Rollback
- Rollback system dotfiles to same state as target dotfiles. 
- Defaults to taking individual files, but can use flag to rollback all files.
- Flag for interactive rollback:
  - Adds extra step where all files are first listed with a numbered index.
  - User provides index number(s) for the files they wish to rollback.
- By default, a folder containing all the original system dotfiles should be created as a backup:
  - Uses backup command functionality for this.
  - This can be disabled with a flag.
  - User can specify backup folder location and extension.

### Restore
- Takes a backup folder and uses it to restore a system to that state.
- Uses rollback command functionality with backup creation disabled.

### Init
- Generates a basic and empty config file.
- Can specify location with index:
  - --help for command will list possible file locations.
  - User passes flag with index for where to generate file (e.g. docon init -l 1).
- Includes docon directory by default. 
- Checks if any config file exists anywhere already, and warns user that new config file with override/be overriden.
- Fails if config file already exists at specified location.

### Verify
- Check config file syntax.
- Notify if excluded files are already present in repo.
- Automatically runs when using any command that requires config file information.

### Pkg
- Generates a packages list in target directory.
- Package list is generated based on Linux distribution being used.
- Add a withHistory flag which, instead of replacing previous file, creates a new one with the date in file name (e.g. pkglist-26.06.21)
- Customise package list location in config file.
- Flag which removes version number from pkglist entries.
- Includes "Added", "Removed", and "Updated" sections when printing differences.
- Flag for quiet mode where no output is printed (just an updated/up-to-date print).
- Includes shell commands for other distros.
- Flag which user can use for passing their own package generation command. This can also be customised in config.yaml.