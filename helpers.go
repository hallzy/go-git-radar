package main

import (
    "fmt"
    "strconv"
    "strings"
    "regexp"
)

// Helper to convert a numeric string to an unsigned int
func str2int(str string) uint {
    ret, err := strconv.ParseUint(str, 10, 32);

    if (err != nil) {
        panic("String [" + str + "] could not be converted to an uint");
    }

    return uint(ret);
}

// Helper to convert an unsigned int to a string
func int2str(num uint) string {
    return fmt.Sprintf("%d", num);
}

// Helper to trim a string
func trim(str string) string {
    return strings.TrimSpace(string(str));
}

// Given how many commits behind or ahead, return the formatted string used in
// the prompt for remote info
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

// Return the formatted prompt string for branch information
func getBranchInfo(remoteBranch string, localBranch string) string {
    if (remoteBranch == "") {
        // If no remote, report the local branch as no upstream
        return fmt.Sprintf(REMOTE_NOT_UPSTREAM, localBranch);
    }

    return fmt.Sprintf(BRANCH_FORMAT, localBranch);
}

// Return the formatted prompt string for local info
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

// Return the formatted prompt string for changed file information (modified,
// deleted etc)
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

// Return the formatted prompt string for conflicting files
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

// Return the formatted prompt string for changed staged files
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

// Return the formatted prompt string for changed unstaged files
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

// Print out the whole prompt given some git Data
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

// Easy to use and remember regex function
func ezRegex(regex string, target string) bool {
    ret, _    := regexp.MatchString(regex, target);
    return ret;
}

// Function to parse the raw git status lines into usable information in a
// GitStatus structure
func parseGitStatus(lines []string) GitStatus {
    ret := GitStatus{};

    for _, line := range lines {
        // STAGED
        if (ezRegex("^M[^M] ", line)) {
            ret.stagedModified += 1;
        }

        if (ezRegex("^A[^A] ", line)) {
            ret.stagedAdded += 1;
        }

        if (ezRegex("^D[^D] ", line)) {
            ret.stagedDeleted += 1;
        }

        if (ezRegex("^R[^R] ", line)) {
            ret.stagedRenamed += 1;
        }

        if (ezRegex("^C[^C] ", line)) {
            ret.stagedCopied += 1;
        }

        if (ezRegex("^T[^T] ", line)) {
            ret.stagedTypeChanged += 1;
        }

        // UNSTAGED

        if (ezRegex("^[^M]M ", line)) {
            ret.unstagedModified += 1;
        }

        if (ezRegex("^[^D]D ", line)) {
            ret.unstagedDeleted += 1;
        }

        if (ezRegex("^[^T]T ", line)) {
            ret.unstagedTypeChanged += 1;
        }

        // CONFLICT

        if (ezRegex("^[^U]U ", line)) {
            ret.conflictUs += 1;
        }

        if (ezRegex("^U[^U] ", line)) {
            ret.conflictThem += 1;
        }

        if (ezRegex("(UU|AA|DD)", line)) {
            ret.conflictBoth += 1;
        }

        // UNTRACKED

        if (ezRegex("^\\?\\? ", line)) {
            ret.untracked += 1;
        }
    }

    return ret;
}
