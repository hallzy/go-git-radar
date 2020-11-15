package main

import (
    "fmt"
    "regexp"
    "strings"
    "unicode/utf8"
)

// Remove colour codes from string
func clean(str string) string {
    reg1, err := regexp.Compile("\x01\033\\[[^m]+");
    if (err != nil) {
        panic("Failed to clean colour codes from string.");
    }

    reg2, err := regexp.Compile("m\x02");
    if (err != nil) {
        panic("Failed to clean colour codes from string.");
    }

    str1 := reg1.ReplaceAllString(str, "");
    return reg2.ReplaceAllString(str1, "");
}

// Find the length of the string (this works for strings with multibyte
// characters)
func strlen(str string) uint {
    return uint(utf8.RuneCountInString(clean(str)));
}

// Add lengths to a list of Examples
func insertLengths(examples []Example) []Example {
    for idx, example := range examples {
        examples[idx].length = strlen(example.prompt);
    }

    return examples;
}

// Find the max length example
func maxLength(examples []Example) uint {
    var max uint = 0;

    for _, example := range examples {
        if (example.length > max) {
            max = example.length;
        }
    }

    return max;
}

type Example struct {
    prompt      string;
    description string;
    length      uint;
}


// This function prints out the help for this tool
func help() {
    examples := []Example {
        {
            description: "Newly created repository. No remote branches.",
            prompt: showPrompt(newGitData(GitData{
                branches: Branches{
                    local: "master",
                },
            })),
        },
        {
            description: "You are in your master branch and are tracking origin master",
            prompt: showPrompt(newGitData(GitData{
                branches: Branches{
                    remote: RemoteBranch{
                        remote: "origin",
                        branch: "master",
                    },
                    local: "master",
                },
            })),
        },
        {
            description: "Created a new branch, but it isn't tracking any remote branches yet.",
            prompt: showPrompt(newGitData(GitData{
                branches: Branches{
                    local: "my-branch",
                },
            })),
        },
        {
            description: "2 New files that aren't being tracked yet. Branch is tracking origin/my-branch and parent is origin/master",
            prompt: showPrompt(newGitData(GitData{
                branches: Branches{
                    remote: RemoteBranch{
                        remote: "origin",
                        branch: "my-branch",
                    },
                    local: "my-branch",
                },
                status: GitStatus{
                    untracked: 2,
                },
            })),
        },
        {
            description: "Same as previous, but now the parent is upstream/dev",
            prompt: showPrompt(newGitData(GitData{
                branches: Branches{
                    remote: RemoteBranch{
                        remote: "origin",
                        branch: "my-branch",
                    },
                    parent: RemoteBranch{
                        remote: "upstream",
                        branch: "dev",
                    },
                    local: "my-branch",
                },
                status: GitStatus{
                    untracked: 2,
                },
            })),
        },
        {
            description: "1 new file staged to commit and 3 modifications that we still need to `git add`",
            prompt: showPrompt(newGitData(GitData{
                branches: Branches{
                    remote: RemoteBranch{
                        remote: "origin",
                        branch: "my-branch",
                    },
                    local: "my-branch",
                },
                status: GitStatus{
                    stagedAdded:      1,
                    unstagedModified: 3,
                },
            })),
        },
        {
            description: "3 Commits waiting to be pushed to remote, origin/my-branch is behind origin/master by 2 commits.",
            prompt: showPrompt(newGitData(GitData{
                branches: Branches{
                    remote: RemoteBranch{
                        remote: "origin",
                        branch: "my-branch",
                    },
                    local: "my-branch",
                },
                remoteBehind: 2,
                localAhead: 3,
            })),
        },
        {
            description: "our commits pushed up, my-branch and its parent have diverged",
            prompt: showPrompt(newGitData(GitData{
                branches: Branches{
                    remote: RemoteBranch{
                        remote: "origin",
                        branch: "my-branch",
                    },
                    local: "my-branch",
                },
                remoteBehind: 2,
                remoteAhead: 3,
            })),
        },
        {
            description: "mid rebase, we are detached and have 3 conflicts caused by US and 2 caused by THEM",
            prompt: showPrompt(newGitData(GitData{
                branches: Branches{
                    local: "detached@94eac67",
                },
                status: GitStatus{
                    conflictThem: 2,
                    conflictUs:   3,
                },
            })),
        },
        {
            description: "rebase complete, our rewritten commits now need pushed up",
            prompt: showPrompt(newGitData(GitData{
                branches: Branches{
                    remote: RemoteBranch{
                        remote: "origin",
                        branch: "my-branch",
                    },
                    local: "my-branch",
                },
                remoteBehind: 2,
                remoteAhead: 3,
                localBehind: 3,
                localAhead: 5,
            })),
        },
        {
            description: "origin/my-branch is 3 commits ahead of it's parent branch",
            prompt: showPrompt(newGitData(GitData{
                branches: Branches{
                    remote: RemoteBranch{
                        remote: "origin",
                        branch: "my-branch",
                    },
                    local: "my-branch",
                },
                remoteAhead: 3,
            })),
        },
        {
            description: "You have 3 stashes stored",
            prompt: showPrompt(newGitData(GitData{
                branches: Branches{
                    remote: RemoteBranch{
                        remote: "origin",
                        branch: "my-branch",
                    },
                    local: "my-branch",
                },
                stash: 3,
            })),
        },
    };

    examplesWithLengths := insertLengths(examples);
    maxLength := maxLength(examplesWithLengths);

    fmt.Println("git-radar - a heads up display for git");
    fmt.Println("");
    fmt.Println("examples:");

    var padding uint;
    for _, example := range examplesWithLengths {
        padding = maxLength - example.length + 2;
        fmt.Println(example.prompt + strings.Repeat(" ", int(padding)) + "# " + example.description);
    }

    fmt.Println("");
    fmt.Println("usage:");
    fmt.Println("  git-radar [help|fetch]");
    fmt.Println("");
    fmt.Println("  fetch  # Fetches your repo asynchronously in the background every 5 mins");
    fmt.Println("  help   # Output this help text.");
    fmt.Println("");
    fmt.Println("Bash example:");
    fmt.Println("  export PS1=\"\\W\\$(git-radar bash fetch) \"");
    fmt.Println("");
    fmt.Println("  This will show your current directory and the full git-radar.");
    fmt.Println("  As an added benefit, if you are in a repo, it will asynchronously");
    fmt.Println("  run `git fetch` every 5 mins, so that you are never out of date.");
    fmt.Println("");
}
