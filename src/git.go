package main

import (
    "strings"
    "os/exec"
)

// =============================================================================
// ALL FUNCTIONS IN THIS FILE PERFORM SIDE EFFECTS
// =============================================================================

// This function populates a GitData structure with current git information
func getGitData() GitData {
    // Get the location of the .git folder
    var dotGit string = dotGit();
    var isRepo bool   = isRepo(dotGit);

    // If we aren't in a repo, then don't bother continuing. Just set that this
    // isn't a repo
    if (!isRepo) {
        return GitData {
            isRepo: isRepo,
            dotGit: dotGit,
        }
    }

    // This is a repo, so set everything else

    // Save these in variables as they are needed for function calls below
    var localBranch  string       = getLocalBranchName();
    var remoteBranch RemoteBranch = getRemoteBranchName(localBranch);
    var parentBranch RemoteBranch = getParentRemote();

    var parentFull = getFullRemote(parentBranch);
    var remoteFull = getFullRemote(remoteBranch);

    return GitData {
        isRepo:     isRepo,
        dotGit:     dotGit,
        branches:   Branches{
            local:  localBranch,
            remote: remoteBranch,
            parent: parentBranch,
        },
        remoteAhead:  howFarAheadRemote(remoteFull, parentFull),
        remoteBehind: howFarBehindRemote(remoteFull, parentFull),
        localAhead:   howFarAheadLocal(remoteFull),
        localBehind:  howFarBehindLocal(remoteFull),
        status:       getGitStatus(),
        stash:        gitStash(),
    };
}

func isRepo(dotGit string) bool {
    // If no dot git path, or we are in the .git folder, then return false
    if (dotGit == "" || dotGit == ".") {
        return false;
    }

    cwd, err := runCmd("pwd");
    if (err != nil) {
        panic("Failed to retrieve the current working directory.");
    }


    // consider this not a repo if we are inside of a .git folder
    return !strings.HasPrefix(cwd, dotGit + "/");
}

// Returns the path to the .git folder for the CWD
func dotGit() string {
    out, err := runCmd("git rev-parse --git-dir");
    if (err != nil) {
        return "";
    }

    return out;
}

// Runs git fetch if enough time has passed
func fetch(dotGit string) bool {
    var now uint = now();

    // How long has it been since the last fetch was made?
    var secSinceLastFetch  uint = now - getLastFetchTime(dotGit);

    // If it hasn't been GIT_RADAR_FETCH_TIME seconds since the last fetch yet,
    // then don't fetch
    if (secSinceLastFetch < GIT_RADAR_FETCH_TIME) {
        return false
    }

    // run the fetch, and save the new fetch time in the fetch time file
    runCmdConcurrent("git fetch --quiet");
    recordNewFetchTime(dotGit);

    return true;
}

// Gets the currently checked out branch name
func getLocalBranchName() string {
    // If no errors, then local branch is the result of the below command
    out1, err := runCmd("git symbolic-ref --short HEAD");
    if (err == nil && out1 != "") {
        return out1;
    }

    // If an error occurred above, then run this command which will give us a
    // short hash
    out2, err := runCmd("git rev-parse --short HEAD");
    if (err != nil) {
        panic("Failed to retrieve branch name: " + err.Error());
    }

    // Return the hash and state detached head
    return "detached@" + out2;
}

// Get the name of the remote tracking branch
func getRemoteBranchName(localBranch string) RemoteBranch {
    var ret RemoteBranch;
    // For the input branch, find out what the remote name is for it (probably
    // origin)
    remote, err := runCmd("git config --get branch." + localBranch + ".remote");
    if (err != nil) {
        return ret;
    }

    // If no remote, then assume no remote branch tracking
    if (remote == "") {
        return ret;
    }

    // Get the branch name that the local branch tracks at the remote we found
    // above
    out, err := runCmd("git config --get branch." + localBranch + ".merge");
    if (err != nil) {
        return ret;
    }

    // The above will prefix the branch with 'refs/heads/'. We don't care about
    // that, so remove it
    var remoteMergeBranch string = strings.Replace(out, "refs/heads/", "", -1);
    if (remoteMergeBranch == "") {
        return ret;
    }

    // Return the remote and branch name together
    ret.remote = remote;
    ret.branch = remoteMergeBranch;

    return ret;
}

// Get the remote branch that the current branch is based on
func getParentRemote() RemoteBranch {
    // Every branch is based on something, so if we can't find what that is,
    // assume we are basing it off of origin master
    defaultRemote := RemoteBranch{
        remote: "origin",
        branch: "master",
    }

    // Get the name of the currently checked out branch
    out1, err := runCmd("git rev-parse --abbrev-ref HEAD");
    if (err != nil) {
        return defaultRemote;
    }

    // Check to see if we have a git-radar tracking config to see if the parent
    // remote branch was saved
    out2, err := runCmd("git config --local branch." + out1 + ".git-radar-tracked-remote");
    if (err != nil) {
        return defaultRemote;
    }

    if (out2 == "") {
        return defaultRemote;
    }

    // Split into 2 strings only, as the branch could have slashes in its name
    remoteBranch := strings.SplitN(out2, "/", 2);

    return RemoteBranch{
        remote: remoteBranch[0],
        branch: remoteBranch[1],
    };
}

// Common function for checking how far ahead/behind remote/local is.
func aheadBehindHelper(toBranch string, fromBranch string, isAhead bool) uint {
    // If any of the branches are empty, or the branches are the same then there
    // is nothing to show.
    if (toBranch == "" || fromBranch == "" || toBranch == fromBranch) {
        return 0;
    }

    var side string;

    // If we are checking how far ahead, then we want to use the --right-only
    // option, otherwise --left-only
    if (isAhead) {
        side = "right-only";
    } else {
        side = "left-only";
    }

    out, err := runCmd("git rev-list --" + side + " --count " + fromBranch + "..." + toBranch);
    if (err != nil) {
        return 0;
    }

    return str2int(out);
}

func howFarAheadRemote(remoteBranch string, parentRemote string) uint {
    return aheadBehindHelper(remoteBranch, parentRemote, true);
}

func howFarBehindRemote(remoteBranch string, parentRemote string) uint {
    return aheadBehindHelper(remoteBranch, parentRemote, false);
}

func howFarAheadLocal(remoteBranch string) uint {
    return aheadBehindHelper("HEAD", remoteBranch, true);
}

func howFarBehindLocal(remoteBranch string) uint {
    return aheadBehindHelper("HEAD", remoteBranch, false);
}

// Get the results of a git status of the current repo and parse the results
// into the GitStatus structure
func getGitStatus() GitStatus {
    porcelain, err := runCmdNoTrim("git status --porcelain");
    if (err != nil) {
        panic("Couldn't get status: " + err.Error());
    }

    if (porcelain == "") {
        return GitStatus{};
    }

    // Split the string on new lines so that we have a list of lines of output
    lines := strings.Split(porcelain, "\n");

    // parse the lines and add everything to a GitStatus structure
    return parseGitStatus(lines);
}

// Run git stash and count how many stashes there are
func gitStash() uint {
    out, err := runCmdNoTrim("git stash list");
    if (err != nil) {
        panic("Failed to get stash info: " + err.Error());
    }

    if (out == "") {
        return 0;
    }

    return countNewLines(out);
}

// Read the fetch time file to see when the last fetch was run
func getLastFetchTime(dotGit string) uint {
    var file string = dotGit + "/git_radar_last_fetch_time";

    if (!fileExists(file)) {
        return 0;
    }

    return str2int(fileRead(file));
}

// Record a new fetch time in the fetch time file
func recordNewFetchTime(dotGit string) bool {
    var file string = dotGit + "/git_radar_last_fetch_time";
    var now  uint    = now();

    fileWrite(file, int2str(now));
    return true;
}

// Run a command, but don't wait for it to finish
func runCmdConcurrent(cmdStr string) {
    var cmdArgs []string = strings.Split(cmdStr, " ");
    cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...);
    err := cmd.Start();

    if (err != nil) {
        panic("Concurrent command [" + cmdStr + "] failed.");
    }
}

// Run a command, but don't trim the output
func runCmdNoTrim(cmdStr string) (string, error) {
    var cmdArgs []string = strings.Split(cmdStr, " ");
    cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...);
    out, err := cmd.Output();

    return string(out), err;
}

// Run a command, trim the output
func runCmd(cmdStr string) (string, error) {
    o, e := runCmdNoTrim(cmdStr);

    return trim(o), e;
}
