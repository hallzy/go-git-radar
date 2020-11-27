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

func getFullRemote(remoteBranch RemoteBranch) string {
    if (remoteBranch.remote == "" || remoteBranch.branch == "") {
        return "";
    }

    return remoteBranch.remote + "/" + remoteBranch.branch;
}

func getFullName(remoteBranch RemoteBranch) string {
    if (remoteBranch.remote == "" || remoteBranch.branch == "") {
        return "";
    }

    if (remoteBranch.remote == "origin") {
        return remoteBranch.branch;
    }

    return remoteBranch.remote + "/" + remoteBranch.branch;
}

// Insert data into a format string
// The keys of placeholderMap are the names of placeholders that are in str
// The values of placeholderMap are what to replace those placeholders with
func insertData(str string, placeholderMap FormatData) string {
    var ret string = str;
    for placeholder, value := range placeholderMap {
        ret = strings.ReplaceAll(ret, "%%" + placeholder + "%%", value);
    }
    return ret;
}

// Given how many commits behind or ahead, return the formatted string used in
// the prompt for remote info
func getRemoteInfo(branches Branches, remoteBehind uint, remoteAhead uint) string {
    var parent string = getFullName(branches.parent);
    var remote string = getFullName(branches.remote);

    // If no remote, report the local branch as no upstream
    if (remote == "") {
        return insertData(REMOTE_NOT_UPSTREAM, FormatData{
            "REMOTE_BRANCH": branches.local,
        });
    }

    if (remoteBehind > 0 && remoteAhead > 0) {
        return insertData(REMOTE_DIVERGED, FormatData{
            "PARENT_REMOTE_BRANCH": parent,
            "REMOTE_BEHIND":        int2str(remoteBehind),
            "REMOTE_AHEAD":         int2str(remoteAhead),
            "REMOTE_BRANCH":        remote,
        });
    }

    if (remoteAhead > 0) {
        return insertData(REMOTE_AHEAD, FormatData{
            "PARENT_REMOTE_BRANCH": parent,
            "REMOTE_AHEAD":         int2str(remoteAhead),
            "REMOTE_BRANCH":        remote,
        });
    }

    if (remoteBehind > 0) {
        return insertData(REMOTE_BEHIND, FormatData{
            "PARENT_REMOTE_BRANCH": parent,
            "REMOTE_BEHIND":        int2str(remoteBehind),
            "REMOTE_BRANCH":        remote,
        });
    }

    if (remote != parent) {
        return insertData(REMOTE_EQUAL, FormatData{
            "PARENT_REMOTE_BRANCH": parent,
            "REMOTE_BRANCH":        remote,
        });
    }

    return insertData(REMOTE_SAME, FormatData{
        "REMOTE_BRANCH": remote,
    });
}

// Return the formatted prompt string for local info
func getLocalInfo(localBehind uint, localAhead uint) string {
    if (localBehind > 0 && localAhead > 0) {
        return insertData(LOCAL_DIVERGED, FormatData{
            "LOCAL_BEHIND": int2str(localBehind),
            "LOCAL_AHEAD": int2str(localAhead),
        });
    }

    if (localAhead > 0) {
        return insertData(LOCAL_AHEAD, FormatData{
            "LOCAL_AHEAD": int2str(localAhead),
        });
    }

    if (localBehind > 0) {
        return insertData(LOCAL_BEHIND, FormatData{
            "LOCAL_BEHIND": int2str(localBehind),
        });
    }

    return "";
}

// Return the formatted prompt string for local info
func showUntracked(gitStatus GitStatus) string {
    if (gitStatus.untracked == 0) {
        return "";
    }

    return UNTRACKED_PREFIX + insertData(CHANGES_UNTRACKED, FormatData{
        "COUNT":  int2str(gitStatus.untracked),
        "SYMBOL": UNTRACKED_SYM,
    }) + UNTRACKED_SUFFIX;
}

// Return the formatted prompt string for conflicting files
func showConflicted(gitStatus GitStatus) string {
    var conflicted string = "";

    if (gitStatus.conflictUs > 0) {
        conflicted += insertData(CHANGES_CONFLICTED, FormatData{
            "COUNT":  int2str(gitStatus.conflictUs),
            "SYMBOL": CONFLICT_US_SYM,
        });
    }

    if (gitStatus.conflictThem > 0) {
        conflicted += insertData(CHANGES_CONFLICTED, FormatData{
            "COUNT":  int2str(gitStatus.conflictThem),
            "SYMBOL": CONFLICT_THEM_SYM,
        });
    }

    if (gitStatus.conflictBoth > 0) {
        conflicted += insertData(CHANGES_CONFLICTED, FormatData{
            "COUNT":  int2str(gitStatus.conflictBoth),
            "SYMBOL": CONFLICT_BOTH_SYM,
        });
    }

    if (conflicted == "") {
        return "";
    }
    return CONFLICTED_PREFIX + conflicted + CONFLICTED_SUFFIX;
}

// Return the formatted prompt string for changed staged files
func showStaged(gitStatus GitStatus) string {
    var staged string = "";

    if (gitStatus.stagedAdded > 0) {
        staged += insertData(CHANGES_STAGED, FormatData{
            "COUNT":  int2str(gitStatus.stagedAdded),
            "SYMBOL": STAGED_ADDED_SYM,
        });
    }

    if (gitStatus.stagedDeleted > 0) {
        staged += insertData(CHANGES_STAGED, FormatData{
            "COUNT":  int2str(gitStatus.stagedDeleted),
            "SYMBOL": STAGED_DELETED_SYM,
        });
    }

    if (gitStatus.stagedModified > 0) {
        staged += insertData(CHANGES_STAGED, FormatData{
            "COUNT":  int2str(gitStatus.stagedModified),
            "SYMBOL": STAGED_MODIFIED_SYM,
        });
    }

    if (gitStatus.stagedRenamed > 0) {
        staged += insertData(CHANGES_STAGED, FormatData{
            "COUNT":  int2str(gitStatus.stagedRenamed),
            "SYMBOL": STAGED_RENAMED_SYM,
        });
    }

    if (gitStatus.stagedCopied > 0) {
        staged += insertData(CHANGES_STAGED, FormatData{
            "COUNT":  int2str(gitStatus.stagedCopied),
            "SYMBOL": STAGED_COPIED_SYM,
        });
    }

    if (gitStatus.stagedTypeChanged > 0) {
        staged += insertData(CHANGES_STAGED, FormatData{
            "COUNT":  int2str(gitStatus.stagedTypeChanged),
            "SYMBOL": STAGED_TYPE_CHANGED_SYM,
        });
    }

    if (staged == "") {
        return "";
    }

    return STAGED_PREFIX + staged + STAGED_SUFFIX;
}

// Return the formatted prompt string for changed unstaged files
func showUnstaged(gitStatus GitStatus) string {
    var unstaged string = "";

    if (gitStatus.unstagedDeleted > 0) {
        unstaged += insertData(CHANGES_UNSTAGED, FormatData{
            "COUNT":  int2str(gitStatus.unstagedDeleted),
            "SYMBOL": UNSTAGED_DELETED_SYM,
        });
    }

    if (gitStatus.unstagedModified > 0) {
        unstaged += insertData(CHANGES_UNSTAGED, FormatData{
            "COUNT":  int2str(gitStatus.unstagedModified),
            "SYMBOL": UNSTAGED_MODIFIED_SYM,
        });
    }

    if (gitStatus.unstagedTypeChanged > 0) {
        unstaged += insertData(CHANGES_UNSTAGED, FormatData{
            "COUNT":  int2str(gitStatus.unstagedTypeChanged),
            "SYMBOL": UNSTAGED_TYPE_CHANGED_SYM,
        });
    }

    if (unstaged == "") {
        return "";
    }
    return UNSTAGED_PREFIX + unstaged + UNSTAGED_SUFFIX;
}

// Print out the whole prompt given some git Data
func showPrompt(git GitData) string {
    remote := getRemoteInfo(git.branches, git.remoteBehind, git.remoteAhead);
    local  := getLocalInfo(git.localBehind, git.localAhead);

    var stash string = "";
    if (git.stash != 0) {
        stash = STASH_PREFIX + insertData(STASH_FORMAT, FormatData{
            "COUNT":  int2str(git.stash),
            "SYMBOL": STASHED_SYM,
        }) + STASH_SUFFIX;
    }

    var fetching string = "";
    if (git.fetching) {
        fetching = FETCH_IN_PROGRESS;
    }

    var staged     string = showStaged(git.status);
    var conflicted string = showConflicted(git.status);
    var unstaged   string = showUnstaged(git.status);
    var untracked  string = showUntracked(git.status);

    return insertData(PROMPT_FORMAT, FormatData{
        "REMOTE_STATUS":      remote,
        "LOCAL_INFO":         local,
        "STASH_STATUS":       stash,
        "STAGED_CHANGES":     staged,
        "CONFLICTED_CHANGES": conflicted,
        "UNSTAGED_CHANGES":   unstaged,
        "UNTRACKED_CHANGES":  untracked,
        "FETCH_IN_PROGRESS":  fetching,
    });
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
        if (ezRegex("^M. ", line)) {
            ret.stagedModified += 1;
        }

        // AA is for conflicts, so can't have 2 A's
        if (ezRegex("^A[^A] ", line)) {
            ret.stagedAdded += 1;
        }

        // DD is for conflicts, so can't have 2 D's
        if (ezRegex("^D[^D] ", line)) {
            ret.stagedDeleted += 1;
        }

        if (ezRegex("^R. ", line)) {
            ret.stagedRenamed += 1;
        }

        if (ezRegex("^C. ", line)) {
            ret.stagedCopied += 1;
        }

        if (ezRegex("^T. ", line)) {
            ret.stagedTypeChanged += 1;
        }

        // UNSTAGED

        if (ezRegex("^.M ", line)) {
            ret.unstagedModified += 1;
        }

        // DD is for conflicts, so can't have 2 D's
        if (ezRegex("^[^D]D ", line)) {
            ret.unstagedDeleted += 1;
        }

        if (ezRegex("^.T ", line)) {
            ret.unstagedTypeChanged += 1;
        }

        // CONFLICT

        if (ezRegex("^[^U]U ", line)) {
            ret.conflictUs += 1;
        }

        if (ezRegex("^U[^U] ", line)) {
            ret.conflictThem += 1;
        }

        if (ezRegex("^(UU|AA|DD)", line)) {
            ret.conflictBoth += 1;
        }

        // UNTRACKED

        if (ezRegex("^\\?\\? ", line)) {
            ret.untracked += 1;
        }
    }

    return ret;
}

// Count the number of lines in a string
func countNewLines(str string) uint {
    return uint(strings.Count(str, "\n"));
}

// Run this function on a created GitData struct to set some defaults
func newGitData(git GitData) GitData {
    if (git.branches.local == "") {
        git.branches.local = "<unset>";
    }

    if (git.branches.parent.remote == "") {
        git.branches.parent.remote = "origin";
    }

    if (git.branches.parent.branch == "") {
        git.branches.parent.branch = "master";
    }

    return git;
}

// Check if cwd is in a git repo
func isRepo(cwd string, dotGit string) bool {
    // If no dot git path, or we are in the .git folder, then return false
    if (dotGit == "" || dotGit == ".") {
        return false;
    }

    // .git is a relative path, which means we are in the root of the repo
    // So this is a repo
    if (dotGit == ".git") {
        return true;
    }

    // Find the root of the repo
    reg1, _ := regexp.Compile("/.git$");
    repoRoot := reg1.ReplaceAllString(dotGit, "");

    // If the cwd isn't inside of the repository root, then we aren't in a repo.
    if (!strings.HasPrefix(cwd, repoRoot)) {
        return false;
    }
    // consider this not a repo if we are inside of a .git folder
    return !strings.HasPrefix(cwd, dotGit + "/");
}
