package main

import (
    "fmt"
    "regexp"
    "strconv"
    "os"
    "os/exec"
    "io/ioutil"
    "time"
    "strings"
)

func getArgs() []string {
    return os.Args[1:];
}

func now() int {
    return int(time.Now().Unix());
}

func str2int(str string) int {
    ret, err := strconv.Atoi(str);

    if (err != nil) {
        panic("String [" + str + "] could not be converted to an int");
    }

    return ret;
}

func int2str(num int) string {
    return strconv.Itoa(num);
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

func colourString(str string, colour string, isBold bool) string {
    var code string;

    switch colour {
        case "black":
            code = ";30";
        case "red":
            code = ";31";
        case "green":
            code = ";32";
        case "yellow":
            code = ";33";
        case "gray":
            fallthrough;
        case "grey":
            code = ";37";
        default:
            code = "";
    }

    var boldness string;
    if (isBold) {
        boldness = "1";
    } else {
        boldness = "0";
    }

    return "\033[" + boldness + code + "m" + str + "\033[0m";
}

func getLastFetchTime() int {
    var file string = dotGit() + "/git_radar_last_fetch_time";

    // Check if file exists
    if (!fileExists(file)) {
        return 0;
    }

    return str2int(fileRead(file));
}

func recordNewFetchTime() bool {
    var file string = dotGit() + "/git_radar_last_fetch_time";
    var now  int    = now();

    fileWrite(file, int2str(now));
    return true;
}

func trim(str string) string {
    return strings.TrimSpace(string(str));
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

func getRemoteInfo() string {
    var remoteBranch string = getRemoteBranchName();

    if (remoteBranch == "") {
        return "";
    }

    var parentRemote string = getParentRemote();

    var behindBy int = howFarBehindRemote(remoteBranch, parentRemote);
    var aheadBy int = howFarAheadRemote(remoteBranch, parentRemote);

    if (behindBy > 0 && aheadBy > 0) {
        return showRemoteDiverged(behindBy, aheadBy);
    }

    if (aheadBy > 0) {
        return showRemoteAhead(aheadBy);
    }

    if (behindBy > 0) {
        return showRemoteBehind(behindBy);
    }

    return "";
}

func getBranchInfo() string {
    var remoteBranch string = getRemoteBranchName();

    if (remoteBranch == "") {
        return fmt.Sprintf(REMOTE_NOT_UPSTREAM, getLocalBranchName());
    }

    return fmt.Sprintf(BRANCH_FORMAT, getLocalBranchName());
}

func getLocalInfo() string {
    var remoteBranch string = getRemoteBranchName();

    if (remoteBranch == "") {
        return "";
    }

    var behindBy int = howFarBehindLocal(remoteBranch);
    var aheadBy int = howFarAheadLocal(remoteBranch);

    if (behindBy > 0 && aheadBy > 0) {
        return showLocalDiverged(behindBy, aheadBy);
    }

    if (aheadBy > 0) {
        return showLocalAhead(aheadBy);
    }

    if (behindBy > 0) {
        return showLocalBehind(behindBy);
    }

    return "";
}

func getStashInfo() string {
    out, err := runCmd("git stash list");
    if (err != nil) {
        panic("Failed to get stash info: " + err.Error());
    }

    if (out == "") {
        return "";
    }

    var count int = strings.Count(out, "\n") + 1;

    return showStash(count);
}

func getChangeInfo() string {
    porcelain, err := runCmdNoTrim("git status --porcelain");
    if (err != nil) {
        panic("Couldn't get status: " + err.Error());
    }

    if (porcelain == "") {
        return "";
    }

    var staged_changes     string = stagedStatus(porcelain);
    var unstaged_changes   string = unstagedStatus(porcelain);
    var conflicted_changes string = conflictedStatus(porcelain);
    var untracked_changes  string = untrackedStatus(porcelain);

    return staged_changes + conflicted_changes + unstaged_changes + untracked_changes;
}

func bool2int(boolean bool) int {
    if (boolean) {
        return 1;
    }
    return 0;
}

func untrackedStatus(gitStatus string) string {
    var untracked int = 0;

    for _, line := range strings.Split(gitStatus, "\n") {
        match1, _    := regexp.MatchString("^\\?\\? ", line);
        untracked += bool2int(match1);
    }

    return showUntracked(untracked);
}

func conflictedStatus(gitStatus string) string {
    var filesUs    int = 0;
    var filesThem  int = 0;
    var filesBoth  int = 0;

    for _, line := range strings.Split(gitStatus, "\n") {
        match1, _    := regexp.MatchString("^[^U]U ", line);
        filesUs += bool2int(match1);

        match2, _     := regexp.MatchString("^U[^U] ", line);
        filesThem += bool2int(match2);

        match3, _ := regexp.MatchString("^(UU|AA|DD) ", line);
        filesBoth += bool2int(match3);
    }

    return showConflicted(filesUs, filesThem, filesBoth);
}

func unstagedStatus(gitStatus string) string {
    var filesModified    int = 0;
    var filesDeleted     int = 0;
    var filesTypeChanged int = 0;

    for _, line := range strings.Split(gitStatus, "\n") {
        match1, _    := regexp.MatchString("^[^M]M ", line);
        filesModified += bool2int(match1);

        match2, _     := regexp.MatchString("^[^D]D ", line);
        filesDeleted += bool2int(match2);

        match3, _ := regexp.MatchString("^[^T]T ", line);
        filesTypeChanged += bool2int(match3);
    }

    return showUnstaged(filesDeleted, filesModified, filesTypeChanged);
}

func stagedStatus(gitStatus string) string {
    var filesModified    int = 0;
    var filesAdded       int = 0;
    var filesDeleted     int = 0;
    var filesRenamed     int = 0;
    var filesCopied      int = 0;
    var filesTypeChanged int = 0;

    for _, line := range strings.Split(gitStatus, "\n") {
        match1, _    := regexp.MatchString("^M[^M] ", line);
        filesModified += bool2int(match1);

        match2, _       := regexp.MatchString("^A[^A] ", line);
        filesAdded += bool2int(match2);

        match3, _     := regexp.MatchString("^D[^D] ", line);
        filesDeleted += bool2int(match3);

        match4, _     := regexp.MatchString("^R[^R] ", line);
        filesRenamed += bool2int(match4);

        match5, _      := regexp.MatchString("^C[^C] ", line);
        filesCopied += bool2int(match5);

        match6, _ := regexp.MatchString("^T[^T] ", line);
        filesTypeChanged += bool2int(match6);
    }

    return showStaged(filesAdded, filesDeleted, filesModified, filesRenamed, filesCopied, filesTypeChanged);
}
