package main

// Custom type to keep track of important Git information for this program
type GitData struct {
    // Whether or not the cwd is a git repository
    isRepo       bool

    // The path to the .git folder of the repo if isRepo is true
    dotGit       string

    // The name of the currently checked out local branch
    localBranch  string

    // The name of the remote branch that the currently checked out local branch
    // is tracking
    remoteBranch string

    // The remote branch that the currently checked out local branch is based
    // off of (the branch that it was branched from)
    parentBranch string

    // How many commits the parent remote branch is ahead of the local
    remoteAhead  uint

    // How many commits the parent remote branch is behind the local
    remoteBehind uint

    // How many commits the remote branch is ahead of the local
    localAhead   uint

    // How many commits the remote branch is behind the local
    localBehind  uint

    // Git status information
    status       GitStatus

    // How many stashes we have
    stash        uint
}

// Custom type to save important git status information
type GitStatus struct {
    // How many untracked files are there
    untracked uint

    // How many files we have that are staged which are added, deleted,
    // modified, etc
    stagedAdded       uint
    stagedDeleted     uint
    stagedModified    uint
    stagedRenamed     uint
    stagedCopied      uint
    stagedTypeChanged uint

    // How many files we have that are NOT staged which are deleted,
    // modified or have had there type changed
    unstagedDeleted     uint
    unstagedModified    uint
    unstagedTypeChanged uint

    // How many files we have that are in conflict
    conflictUs   uint
    conflictThem uint
    conflictBoth uint
}
