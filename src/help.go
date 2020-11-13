package main

import (
    "fmt"
)

// This function prints out the help for this tool
func help() {
    prompts := []string {
        showPrompt(GitData{
            remoteBranch: "origin/master",
            localBranch: "master",
        }),
        showPrompt(GitData{
            localBranch: "my-branch",
        }),
        showPrompt(GitData{
            remoteBranch: "origin/my-branch",
            localBranch: "my-branch",
            status: GitStatus{
                untracked: 2,
            },
        }),
        showPrompt(GitData{
            remoteBranch: "origin/my-branch",
            localBranch: "my-branch",
            status: GitStatus{
                stagedAdded:      1,
                unstagedModified: 3,
            },
        }),
        showPrompt(GitData{
            remoteBranch: "origin/my-branch",
            localBranch: "my-branch",
            remoteBehind: 2,
            localAhead: 3,
        }),
        showPrompt(GitData{
            remoteBranch: "origin/my-branch",
            localBranch: "my-branch",
            remoteBehind: 2,
            remoteAhead: 3,
        }),
        showPrompt(GitData{
            localBranch: "detached@94eac67",
            status: GitStatus{
                conflictThem: 2,
                conflictUs:   3,
            },
        }),
        showPrompt(GitData{
            remoteBranch: "origin/my-branch",
            localBranch: "my-branch",
            remoteBehind: 2,
            remoteAhead: 3,
            localBehind: 3,
            localAhead: 5,
        }),
        showPrompt(GitData{
            remoteBranch: "origin/my-branch",
            localBranch: "my-branch",
            remoteAhead: 3,
        }),
        showPrompt(GitData{
            remoteBranch: "origin/my-branch",
            localBranch: "my-branch",
            stash: 3,
        }),
    };

    fmt.Println("git-radar - a heads up display for git");
    fmt.Println("");
    fmt.Println("examples:");

    fmt.Println(prompts[0], "\t\t\t\t# You are on the master branch and everything is clean");
    fmt.Println(prompts[1], "\t\t\t# Fresh branch that we haven't pushed upstream");
    fmt.Println(prompts[2], "\t\t\t\t# Two files created that aren't tracked by git");
    fmt.Println(prompts[3], "\t\t\t# 1 new file staged to commit and 3 modifications that we still need to `git add`");
    fmt.Println(prompts[4], "\t\t\t# 3 commits made locally ready to push up while master is ahead of us by 2");
    fmt.Println(prompts[5], "\t\t\t# our commits pushed up, master and my-branch have diverged");
    fmt.Println(prompts[6], "\t# mid rebase, we are detached and have 3 conflicts caused by US and 2 caused by THEM");
    fmt.Println(prompts[7], "\t\t\t# rebase complete, our rewritten commits now need pushed up");
    fmt.Println(prompts[8], "\t\t\t# origin/my-branch is up to date with master and has our 3 commits waiting merge");
    fmt.Println(prompts[9], "\t\t\t\t# You have 3 stashes stored");

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
