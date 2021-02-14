package main


// TODO: IMPROVE THIS FILE

import(
    "testing"
    "strings"
    "os"
    "fmt"
)

func TestDotGit(T *testing.T) {
    var cwd string = getCwd();

    defer os.Chdir(cwd);

    var dotGitStr string = dotGit();

    if (!strings.HasSuffix(dotGitStr, ".git")) {
        T.Errorf("dotGit(): Failed to get .git folder path");
        return;
    }

    // Assume that this is not a git repo...
    os.Chdir("../../");

    dotGitStr = dotGit();
    if (dotGitStr != "") {
        T.Errorf("dotGit(): Expected no .git folder (should be outside of git repo).");
        return;
    }
}


func TestIsRepo(T *testing.T) {
    if (!isRepo()) {
        T.Errorf("isRepo(): expected true");
    }

    var cwd string = getCwd();
    defer os.Chdir(cwd);

    // Assume that this is not a git repo...
    os.Chdir("../../");

    if (isRepo()) {
        T.Errorf("isRepo(): expected false (shouldn't be in a repo)");
    }
}

// TODO: This is obviously not a good test yet
func TestGetGitData(T *testing.T) {
    // Come back to current directory when done test
    defer os.Chdir(getCwd());

    var data GitData = getGitData();

    if (!data.isRepo) {
        T.Errorf("isRepo(): expected true");
    }

    if (!strings.HasSuffix(data.dotGit, ".git")) {
        T.Errorf("dotGit(): Failed to get .git folder path");
        return;
    }

    // These Values can't really be checked to be correct. The only way to
    // compare them with the current values is to compare them against same
    // functions that gave us these values.
    // The Repo State is also not deterministic so I can't know what these
    // values should be. Just check them yourself.
    fmt.Println("MANUALLY CHECK THESE VALUES WITH THE CURRENT STATE OF YOUR REPOSITORY");
    fmt.Println("=====================================================================");
    fmt.Println("Local Branch:        ", data.branches.local);
    fmt.Println("Remote Branch:       ", data.branches.remote.remote, "/", data.branches.remote.branch, "(exists? ", data.branches.remote.exists, ")");
    fmt.Println("Parent Branch:       ", data.branches.parent.remote, "/", data.branches.parent.branch, "(exists? ", data.branches.parent.exists, ")");
    fmt.Println("Remote Ahead:        ", data.remoteAhead);
    fmt.Println("Remote Behind:       ", data.remoteBehind);
    fmt.Println("Local Ahead:         ", data.localAhead);
    fmt.Println("Local Behind:        ", data.localBehind);
    fmt.Println("Fetching:            ", data.fetching);
    fmt.Println("Status:");
    fmt.Println("    Untracked:       ", data.status.untracked);
    fmt.Println("    Stash:           ", data.stash);
    fmt.Println("    Staged:");
    fmt.Println("        Added:       ", data.status.stagedAdded);
    fmt.Println("        Deleted:     ", data.status.stagedDeleted);
    fmt.Println("        Modified:    ", data.status.stagedModified);
    fmt.Println("        Renamed:     ", data.status.stagedRenamed);
    fmt.Println("        Copied:      ", data.status.stagedCopied);
    fmt.Println("        Type Change: ", data.status.stagedTypeChanged);
    fmt.Println("    Unstaged:");
    fmt.Println("        Deleted:     ", data.status.unstagedDeleted);
    fmt.Println("        Modified:    ", data.status.unstagedModified);
    fmt.Println("        Type Change: ", data.status.unstagedTypeChanged);
    fmt.Println("    Conflict:");
    fmt.Println("        Us:          ", data.status.conflictUs);
    fmt.Println("        Them:        ", data.status.conflictThem);
    fmt.Println("        Both:        ", data.status.conflictBoth);
    fmt.Println("=====================================================================");
}
