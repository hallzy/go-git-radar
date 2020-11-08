package main

import (
    "fmt"
    "strconv"
    "os"
    "io/ioutil"
    "time"
    "strings"
)

func getArgs() []string {
    return os.Args[1:];
}

func now() uint {
    return uint(time.Now().Unix());
}

func str2int(str string) uint {
    ret, err := strconv.ParseUint(str, 10, 32);

    if (err != nil) {
        panic("String [" + str + "] could not be converted to an uint");
    }

    return uint(ret);
}

func int2str(num uint) string {
    return fmt.Sprintf("%d", num);
}

func fileExists(file string) bool {
    _, err := os.Stat(file);

    return ! os.IsNotExist(err);
}

func fileRead(file string) string {
    ret, err := ioutil.ReadFile(file);

    if (err != nil) {
        panic(err.Error());
    }

    return string(ret);
}

func fileWrite(file string, data string) bool {
    err := ioutil.WriteFile(file, []byte(data), 0664);

    if (err != nil) {
        panic(err.Error());
    }

    return true;
}

func trim(str string) string {
    return strings.TrimSpace(string(str));
}

func getRemoteInfo(remoteBehind uint, remoteAhead uint) string {
    if (remoteBehind > 0 && remoteAhead > 0) {
        return fmt.Sprintf(REMOTE_DIVERGED, remoteBehind, remoteAhead);
    }

    if (remoteAhead > 0) {
        return fmt.Sprintf(REMOTE_AHEAD, remoteAhead);
    }

    if (remoteBehind > 0) {
        return fmt.Sprintf(REMOTE_BEHIND, remoteBehind);
    }

    return "";
}

func getBranchInfo(remoteBranch string, localBranch string) string {
    if (remoteBranch == "") {
        return fmt.Sprintf(REMOTE_NOT_UPSTREAM, localBranch);
    }

    return fmt.Sprintf(BRANCH_FORMAT, localBranch);
}

func getLocalInfo(localBehind uint, localAhead uint) string {
    if (localBehind > 0 && localAhead > 0) {
        return fmt.Sprintf(LOCAL_DIVERGED, localBehind, localAhead);
    }

    if (localAhead > 0) {
        return fmt.Sprintf(LOCAL_AHEAD, localAhead);
    }

    if (localBehind > 0) {
        return fmt.Sprintf(LOCAL_BEHIND, localBehind);
    }

    return "";
}

func getChangeInfo(gitStatus GitStatus) string {
    var staged     string = showStaged(gitStatus);
    var unstaged   string = showUnstaged(gitStatus);
    var conflicted string = showConflicted(gitStatus);
    var untracked  string;


    untracked = "";
    if (gitStatus.untracked != 0) {
        untracked = " " + fmt.Sprintf(CHANGES_UNTRACKED, gitStatus.untracked, "A");
    }

    return staged + conflicted + unstaged + untracked;
}

func showConflicted(gitStatus GitStatus) string {
    var conflicted string = "";

    if (gitStatus.conflictUs > 0) {
        conflicted += fmt.Sprintf(CHANGES_CONFLICTED, gitStatus.conflictUs, "U");
    }

    if (gitStatus.conflictThem > 0) {
        conflicted += fmt.Sprintf(CHANGES_CONFLICTED, gitStatus.conflictThem, "T");
    }

    if (gitStatus.conflictBoth > 0) {
        conflicted += fmt.Sprintf(CHANGES_CONFLICTED, gitStatus.conflictBoth, "B");
    }

    if (conflicted == "") {
        return "";
    }
    return " " + conflicted;
}

func showStaged(gitStatus GitStatus) string {
    var staged string = "";

    if (gitStatus.stagedAdded > 0) {
        staged += fmt.Sprintf(CHANGES_STAGED, gitStatus.stagedAdded, "A");
    }

    if (gitStatus.stagedDeleted > 0) {
        staged += fmt.Sprintf(CHANGES_STAGED, gitStatus.stagedDeleted, "D");
    }

    if (gitStatus.stagedModified > 0) {
        staged += fmt.Sprintf(CHANGES_STAGED, gitStatus.stagedModified, "M");
    }

    if (gitStatus.stagedRenamed > 0) {
        staged += fmt.Sprintf(CHANGES_STAGED, gitStatus.stagedRenamed, "R");
    }

    if (gitStatus.stagedCopied > 0) {
        staged += fmt.Sprintf(CHANGES_STAGED, gitStatus.stagedCopied, "C");
    }

    if (gitStatus.stagedTypeChanged > 0) {
        staged += fmt.Sprintf(CHANGES_STAGED, gitStatus.stagedTypeChanged, "TC");
    }

    if (staged == "") {
        return "";
    }

    return " " + staged;
}

func showUnstaged(gitStatus GitStatus) string {
    var unstaged string = "";

    if (gitStatus.unstagedDeleted > 0) {
        unstaged += fmt.Sprintf(CHANGES_UNSTAGED, gitStatus.unstagedDeleted, "D");
    }

    if (gitStatus.unstagedModified > 0) {
        unstaged += fmt.Sprintf(CHANGES_UNSTAGED, gitStatus.unstagedModified, "M");
    }

    if (gitStatus.unstagedTypeChanged > 0) {
        unstaged += fmt.Sprintf(CHANGES_UNSTAGED, gitStatus.unstagedTypeChanged, "TC");
    }

    if (unstaged == "") {
        return "";
    }
    return " " + unstaged;
}

func showPrompt(git GitData) string {
    remote := getRemoteInfo(git.remoteBehind, git.remoteAhead);
    branch := getBranchInfo(git.remoteBranch, git.localBranch);
    local  := getLocalInfo(git.localBehind, git.localAhead);

    var stash string;

    stash = "";
    if (git.stash != 0) {
        stash = fmt.Sprintf(STASH_FORMAT, git.stash);
    }

    change := getChangeInfo(git.status);

    return fmt.Sprintf(PROMPT_FORMAT, PREFIX, remote, branch, local, stash, change, SUFFIX);
}
