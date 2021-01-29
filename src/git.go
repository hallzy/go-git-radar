package main

import (
    "strings"
    "sync"
)

// =============================================================================
// ALL FUNCTIONS IN THIS FILE PERFORM SIDE EFFECTS
// =============================================================================

// This function populates a GitData structure with current git information
func getGitData() GitData {
    var isRepo bool   = isRepo();

    // If we aren't in a repo, then don't bother continuing. Just set that this
    // isn't a repo
    if (!isRepo) {
        return GitData {
            isRepo: isRepo,
        }
    }
    // This is a repo, so set everything else

    var wg sync.WaitGroup;

    // Save these in variables as they are needed for function calls below
    var localAhead   uint;
    var localBehind  uint;
    var localBranch  string;
    var parentBranch RemoteBranch;
    var parentFull   string;
    var remoteAhead  uint;
    var remoteBehind uint;
    var remoteBranch RemoteBranch;
    var remoteFull   string;
    var stash        uint;
    var status       GitStatus;

    // Git calls can be easily the slowest part of this program especially in
    // large repos. So we will try to run them all asynchronously to try and
    // speed this up.
    wg.Add(7);

    go func() {
        defer wg.Done();
        status = getGitStatus();
    }();

    go func() {
        defer wg.Done();
        stash = gitStash();
    }();

    go func() {
        defer wg.Done();

        var wg2 sync.WaitGroup;
        wg2.Add(2);
        go func() {
            defer wg2.Done();
            localBranch         = getLocalBranchName();
            remoteBranch        = getRemoteBranchName(localBranch);
            remoteFull          = getFullRemote(remoteBranch);
            remoteBranch.exists = doesRemoteExist(remoteFull);
        }();

        go func() {
            defer wg2.Done();
            parentBranch        = getParentRemote();
            parentFull          = getFullRemote(parentBranch);
            parentBranch.exists = doesRemoteExist(parentFull);
        }();

        wg2.Wait();

        go func() {
            defer wg.Done();
            remoteAhead = howFarAheadRemote(remoteFull, parentFull);
        }();

        go func() {
            defer wg.Done();
            remoteBehind = howFarBehindRemote(remoteFull, parentFull);
        }();

        go func() {
            defer wg.Done();
            localAhead = howFarAheadLocal(remoteFull);
        }();

        go func() {
            defer wg.Done();
            localBehind = howFarBehindLocal(remoteFull);
        }();
    }();

    wg.Wait();

    var dotGit string = dotGit();

    return GitData {
        isRepo:     isRepo,
        dotGit:     dotGit,
        branches:   Branches{
            local:  localBranch,
            remote: remoteBranch,
            parent: parentBranch,
        },
        remoteAhead:  remoteAhead,
        remoteBehind: remoteBehind,
        localAhead:   localAhead,
        localBehind:  localBehind,
        status:       status,
        stash:        stash,
    };
}

// Returns the path to the .git folder for the CWD
func dotGit() string {
    out, err := runCmd("git rev-parse --absolute-git-dir");
    if (err != nil) {
        return "";
    }

    return out;
}

// Runs git fetch if enough time has passed. Return true if we are fetching, or
// even if we are currently fetching from a previous fetch command
func fetch(dotGit string) bool {
    // If we are already fetching, then return true and no need to go any
    // further
    var fetchingFile string = dotGit + "/radar_fetching";
    if (fileExists(fetchingFile)) {
        return true;
    }

    var now uint = now();

    // How long has it been since the last fetch was made?
    var secSinceLastFetch  uint = now - getLastFetchTime(dotGit);

    // If it hasn't been GIT_RADAR_FETCH_TIME seconds since the last fetch yet,
    // then don't fetch
    if (secSinceLastFetch < GIT_RADAR_FETCH_TIME) {
        return false
    }

    var preFetch   string = "( " + PRE_FETCH_CMD + " )";
    var fetch      string = "( git fetch --quiet && (" + FETCH_SUCCEEDED_CMD + ") || (" + FETCH_FAILED_CMD + ") )";
    var outerFetch string = "( touch " + fetchingFile + " && " + fetch + " || ( " + FETCH_FAILED_CMD + " ) )";

    // run the fetch, and save the new fetch time in the fetch time file
    // We are also creating and removing a file to tell us if we are currently
    // fetching already
    runCmdConcurrent("(" + preFetch + " && " + outerFetch + " || ( " + PRE_FETCH_CMD_FAILED + " )); rm " + fetchingFile);
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
    // hash
    hash, err := runCmd("git rev-parse HEAD");
    if (err != nil) {
        panic("Failed to retrieve branch name: " + err.Error());
    }

    // Check to see if we currently have a tag checked out
    tag, err := runCmd("git describe --tag " + hash + " --exact-match");

    var tagOrHash string;
    if (err == nil) {
        tagOrHash = tag;
    } else {
        tagOrHash = hash[0:6];
    }

    // Return the first part of the hash and state detached head
    return "detached@" + tagOrHash;
}

// Get the name of the remote tracking branch
func getRemoteBranchName(localBranch string) RemoteBranch {
    var ret RemoteBranch;
    // For the input branch, find out what the remote name is for it (probably
    // origin)
    remote, err := runCmd("git config --get branch." + localBranch + ".remote");
    if (err != nil || remote == "") {
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

func doesRemoteExist(remoteFull string) bool {
    if remoteFull == "" {
        return false;
    }

    _, err := runCmd("git show-branch " + remoteFull);
    if (err != nil) {
        return false;
    }
    return true;
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
    if (err != nil || out2 == "") {
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
func aheadBehindHelper(toBranch string, fromBranch string, side string) uint {
    // If any of the branches are empty, or the branches are the same then there
    // is nothing to show.
    if (toBranch == "" || fromBranch == "" || toBranch == fromBranch) {
        return 0;
    }

    out, err := runCmd("git rev-list --" + side + " --count " + fromBranch + "..." + toBranch);
    if (err != nil) {
        return 0;
    }

    return str2int(out);
}

func howFarAheadRemote(remoteBranch string, parentRemote string) uint {
    return aheadBehindHelper(remoteBranch, parentRemote, "right-only");
}

func howFarBehindRemote(remoteBranch string, parentRemote string) uint {
    return aheadBehindHelper(remoteBranch, parentRemote, "left-only");
}

func howFarAheadLocal(remoteBranch string) uint {
    return aheadBehindHelper("HEAD", remoteBranch, "right-only");
}

func howFarBehindLocal(remoteBranch string) uint {
    return aheadBehindHelper("HEAD", remoteBranch, "left-only");
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
    // parse the lines and add everything to a GitStatus structure
    return parseGitStatus(strings.Split(porcelain, "\n"));
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

// Check if cwd is in a git repo
func isRepo() bool {
    out, err := runCmd("git rev-parse --is-inside-work-tree");
    if (err != nil) {
        return false;
    }

    return out == "true";
}
