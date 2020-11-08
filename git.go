package main

import (
    "strings"
    "os/exec"
    "regexp"
)

type GitData struct {
    isRepo       bool
    dotGit       string
    localBranch  string
    remoteBranch string
    parentBranch string
    remoteAhead  uint
    remoteBehind uint
    localAhead   uint
    localBehind  uint
    status       GitStatus
    stash        uint
}

type GitStatus struct {
    untracked uint

    stagedAdded       uint
    stagedDeleted     uint
    stagedModified    uint
    stagedRenamed     uint
    stagedCopied      uint
    stagedTypeChanged uint

    unstagedDeleted     uint
    unstagedModified    uint
    unstagedTypeChanged uint

    conflictUs   uint
    conflictThem uint
    conflictBoth uint
}

func getGitData() GitData {
    var dotGit string = dotGit();
    var isRepo bool   = dotGit != "";

    if (!isRepo) {
        return GitData {
            isRepo: isRepo,
            dotGit: dotGit,
        }
    }

    var localBranch  string = getLocalBranchName();
    var remoteBranch string = getRemoteBranchName(localBranch);
    var parentBranch string = getParentRemote();


    return GitData {
        isRepo:       isRepo,
        dotGit:       dotGit,
        localBranch:  localBranch,
        remoteBranch: remoteBranch,
        parentBranch: parentBranch,
        remoteAhead:  howFarAheadRemote(remoteBranch, parentBranch),
        remoteBehind: howFarBehindRemote(remoteBranch, parentBranch),
        localAhead:   howFarAheadLocal(remoteBranch),
        localBehind:  howFarBehindLocal(remoteBranch),
        status:       getGitStatus(),
        stash:        gitStash(),
    };
}

// Returns whether or not this is a repo, and the location of the .git folder
func dotGit() string {
    out, err := runCmd("git rev-parse --git-dir");
    if (err != nil) {
        return "";
    }

    return out;
}

// returns true if a fetch was made, and false otherwise
func fetch(dotGit string) bool {
    var now                uint = now();
    var secSinceLastFetch  uint = now - getLastFetchTime(dotGit);

    if (secSinceLastFetch < GIT_RADAR_FETCH_TIME) {
        return false
    }

    runCmdConcurrent("git fetch --quiet");

    recordNewFetchTime(dotGit);
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

func getRemoteBranchName(localBranch string) string {
    remote, err := runCmd("git config --get branch." + localBranch + ".remote");
    if (err != nil) {
        return "";
    }

    if (remote == "") {
        return "";
    }

    out, err := runCmd("git config --get branch." + localBranch + ".merge");
    if (err != nil) {
        return "";
    }

    // Remove refs/heads/
    var remoteMergeBranch string = strings.Replace(out, "refs/heads/", "", -1);
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

func howFarAheadRemote(remoteBranch string, parentRemote string) uint {
    if (remoteBranch == "" || parentRemote == "" || remoteBranch == parentRemote) {
        return 0;
    }

    out, err := runCmd("git rev-list --right-only --count " + parentRemote + "..." + remoteBranch);
    if (err != nil) {
        return 0;
    }

    return str2int(out);
}

func howFarBehindRemote(remoteBranch string, parentRemote string) uint {
    if (remoteBranch == "" || parentRemote == "" || remoteBranch == parentRemote) {
        return 0;
    }

    out, err := runCmd("git rev-list --left-only --count " + parentRemote + "..." + remoteBranch);
    if (err != nil) {
        return 0;
    }

    return str2int(out);
}

func howFarAheadLocal(remoteBranch string) uint {
    if (remoteBranch == "") {
        return 0;
    }

    out, err := runCmd("git rev-list --right-only --count " + remoteBranch + "...HEAD");
    if (err != nil) {
        return 0;
    }

    return str2int(out);
}

func howFarBehindLocal(remoteBranch string) uint {
    if (remoteBranch == "") {
        return 0;
    }

    out, err := runCmd("git rev-list --left-only --count " + remoteBranch + "...HEAD");
    if (err != nil) {
        return 0;
    }

    return str2int(out);
}

func ezRegex(regex string, target string) bool {
    ret, _    := regexp.MatchString(regex, target);
    return ret;
}

func getGitStatus() GitStatus {
    porcelain, err := runCmdNoTrim("git status --porcelain");
    if (err != nil) {
        panic("Couldn't get status: " + err.Error());
    }

    if (porcelain == "") {
        return GitStatus{};
    }

    ret := GitStatus{};
    lines := strings.Split(porcelain, "\n");

    for _, line := range lines {
        // STAGED
        if (ezRegex("^M[^M] ", line)) {
            ret.stagedModified += 1;
        } else if (ezRegex("^A[^A] ", line)) {
            ret.stagedAdded += 1;
        } else if (ezRegex("^D[^D] ", line)) {
            ret.stagedDeleted += 1;
        } else if (ezRegex("^R[^R] ", line)) {
            ret.stagedRenamed += 1;
        } else if (ezRegex("^C[^C] ", line)) {
            ret.stagedCopied += 1;
        } else if (ezRegex("^T[^T] ", line)) {
            ret.stagedTypeChanged += 1;
        // UNSTAGED
        } else if (ezRegex("^[^M]M ", line)) {
            ret.unstagedModified += 1;
        } else if (ezRegex("^[^D]D ", line)) {
            ret.unstagedDeleted += 1;
        } else if (ezRegex("^[^T]T ", line)) {
            ret.unstagedTypeChanged += 1;
        // CONFLICT
        } else if (ezRegex("^[^U]U ", line)) {
            ret.conflictUs += 1;
        } else if (ezRegex("^U[^U] ", line)) {
            ret.conflictThem += 1;
        } else if (ezRegex("(UU|AA|DD)", line)) {
            ret.conflictBoth += 1;
        // UNTRACKED
        } else if (ezRegex("^\\?\\? ", line)) {
            ret.untracked += 1;
        }
    }

    return ret;
}

func gitStash() uint {
    out, err := runCmdNoTrim("git stash list");
    if (err != nil) {
        panic("Failed to get stash info: " + err.Error());
    }

    if (out == "") {
        return 0;
    }

    return uint(strings.Count(out, "\n"));
}

func getLastFetchTime(dotGit string) uint {
    var file string = dotGit + "/git_radar_last_fetch_time";

    // Check if file exists
    if (!fileExists(file)) {
        return 0;
    }

    return str2int(fileRead(file));
}

func recordNewFetchTime(dotGit string) bool {
    var file string = dotGit + "/git_radar_last_fetch_time";
    var now  uint    = now();

    fileWrite(file, int2str(now));
    return true;
}

func runCmdConcurrent(cmdStr string) {
    var cmdArgs []string = strings.Split(cmdStr, " ");
    cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...);
    err := cmd.Start();

    if (err != nil) {
        panic("Concurrent command [" + cmdStr + "] failed.");
    }
}

func runCmdNoTrim(cmdStr string) (string, error) {
    var cmdArgs []string = strings.Split(cmdStr, " ");
    cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...);
    out, err := cmd.Output();

    return string(out), err;
}

func runCmd(cmdStr string) (string, error) {
    o, e := runCmdNoTrim(cmdStr);

    return trim(o), e;
}

