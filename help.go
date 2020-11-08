package main

import (
    "fmt"
)
func help() {
    fmt.Println("git-radar - a heads up display for git");
    fmt.Println("");
    fmt.Println("examples:");

    fmt.Print(showPrompt(GitData{
        remoteBranch: "origin/master",
        localBranch: "master",
    }));
    fmt.Println("                            # You are on the master branch and everything is clean");

    fmt.Print(showPrompt(GitData{
        localBranch: "my-branch",
    }));
    fmt.Println("              # Fresh branch that we haven't pushed upstream");

    fmt.Print(showPrompt(GitData{
        remoteBranch: "origin/my-branch",
        localBranch: "my-branch",
        status: GitStatus{
            untracked: 2,
        },
    }));
    fmt.Println("                      # Two files created that aren't tracked by git");

    fmt.Print(showPrompt(GitData{
        remoteBranch: "origin/my-branch",
        localBranch: "my-branch",
        status: GitStatus{
            stagedAdded:      1,
            unstagedModified: 3,
        },
    }));
    fmt.Println("                   # 1 new file staged to commit and 3 modifications that we still need to `git add`");

    fmt.Print(showPrompt(GitData{
        remoteBranch: "origin/my-branch",
        localBranch: "my-branch",
        remoteBehind: 2,
        localAhead: 3,
    }));
    fmt.Println("                # 3 commits made locally ready to push up while master is ahead of us by 2");

    fmt.Print(showPrompt(GitData{
        remoteBranch: "origin/my-branch",
        localBranch: "my-branch",
        remoteBehind: 2,
        remoteAhead: 3,
    }));
    fmt.Println("                 # our commits pushed up, master and my-branch have diverged");

    fmt.Print(showPrompt(GitData{
        localBranch: "detached@94eac67",
        status: GitStatus{
            conflictThem: 2,
            conflictUs:   3,
        },
    }));
    fmt.Println("  # mid rebase, we are detached and have 3 conflicts caused by US and 2 caused by THEM");

    fmt.Print(showPrompt(GitData{
        remoteBranch: "origin/my-branch",
        localBranch: "my-branch",
        remoteBehind: 2,
        remoteAhead: 3,
        localBehind: 3,
        localAhead: 5,
    }));
    fmt.Println("             # rebase complete, our rewritten commits now need pushed up");

    fmt.Print(showPrompt(GitData{
        remoteBranch: "origin/my-branch",
        localBranch: "my-branch",
        remoteAhead: 3,
    }));
    fmt.Println("                   # origin/my-branch is up to date with master and has our 3 commits waiting merge");

    fmt.Print(showPrompt(GitData{
        remoteBranch: "origin/my-branch",
        localBranch: "my-branch",
        stash: 3,
    }));
    fmt.Println("                      # You have 3 stashes stored");

    fmt.Println("");
    fmt.Println("usage:");
    fmt.Println("  git-radar [bash] [fetch]");
    fmt.Println("");
    fmt.Println("  fetch  # Fetches your repo asynchronously in the background every 5 mins");
    fmt.Println("  bash   # Output prompt using Bash style color characters");
    fmt.Println("");
    fmt.Println("Bash example:");
    fmt.Println("  export PS1=\"\\W\\$(git-radar bash fetch) \"");
    fmt.Println("");
    fmt.Println("  This will show your current directory and the full git-radar.");
    fmt.Println("  As an added benefit, if you are in a repo, it will asynchronously");
    fmt.Println("  run `git fetch` every 5 mins, so that you are never out of date.");
    fmt.Println("");
}
