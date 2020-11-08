package main

import (
    "strings"
)

// Returns whether or not this is a repo, and the location of the .git folder
func dotGit() string {
    out, err := runCmd("git rev-parse --git-dir");
    if (err != nil) {
        return "";
    }

    return out;
}

func isRepo() bool {
    return dotGit() != "";
}

// returns true if a fetch was made, and false otherwise
func fetch() bool {
    var now                int = now();
    var secSinceLastFetch  int = now - getLastFetchTime();

    if (secSinceLastFetch < GIT_RADAR_FETCH_TIME) {
        return false
    }

    // TODO: Would be nice to background this so that git-radar can exit while
    // this runs in the background
    runCmdConcurrent("git fetch --quiet");

    recordNewFetchTime();
    return true;
}

func getLocalBranchName() string {
    out1, err := runCmd("git symbolic-ref --short HEAD");
    if (err == nil) {
        return out1;
    }

    out2, err := runCmd("git rev-parse --short HEAD");
    if (err != nil) {
        panic("Failed to retrieve branch name: " + err.Error());
    }
    return "detached@" + out2;
}

func getRemoteName() string {
    out, err := runCmd("git config --get branch." + getLocalBranchName() + ".remote");
    if (err != nil) {
        return "";
    }

    return out;
}

func getRemoteMergeBranch() string {
    out, err := runCmd("git config --get branch." + getLocalBranchName() + ".merge");
    if (err != nil) {
        return "";
    }

    // Remove refs/heads/
    return strings.Replace(out, "refs/heads/", "", -1);
}


func getRemoteBranchName() string {
    var remote string = getRemoteName();
    if (remote == "") {
        return "";
    }

    var remoteMergeBranch string = getRemoteMergeBranch();
    if (remoteMergeBranch == "") {
        return "";
    }

    return remote + "/" + remoteMergeBranch;
}

// Get the remote branch that the current branch is based on
func getParentRemote() string {
    var defaultRemote string = "origin/master";

    out1, err := runCmd("git rev-parse --abbrev-ref HEAD");
    if (err != nil) {
        return defaultRemote;
    }

    out2, err := runCmd("git config --local branch." + out1 + ".git-radar-tracked-remote");
    if (err != nil) {
        return defaultRemote;
    }

    if (out2 == "") {
        return defaultRemote;
    }

    return out2;
}

func howFarAheadRemote(remoteBranch string, parentRemote string) int {
    if (remoteBranch == parentRemote) {
        return 0;
    }

    out, err := runCmd("git rev-list --right-only --count " + parentRemote + "..." + remoteBranch);
    if (err != nil) {
        return 0;
    }

    return str2int(out);
}

func howFarBehindRemote(remoteBranch string, parentRemote string) int {
    if (remoteBranch == parentRemote) {
        return 0;
    }

    out, err := runCmd("git rev-list --left-only --count " + parentRemote + "..." + remoteBranch);
    if (err != nil) {
        return 0;
    }

    return str2int(out);
}

func howFarAheadLocal(remoteBranch string) int {
    out, err := runCmd("git rev-list --right-only --count " + remoteBranch + "...HEAD");
    if (err != nil) {
        return 0;
    }

    return str2int(out);
}

func howFarBehindLocal(remoteBranch string) int {
    out, err := runCmd("git rev-list --left-only --count " + remoteBranch + "...HEAD");
    if (err != nil) {
        return 0;
    }

    return str2int(out);
}
