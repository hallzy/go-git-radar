package main

import (
    "regexp"
    "strings"
    "unicode/utf8"
)

// Remove colour codes from string
func clean(str string) string {
    reg1, _ := regexp.Compile("\x01\033\\[[^m]+");
    reg2, _ := regexp.Compile("m\x02");
    return reg2.ReplaceAllString(reg1.ReplaceAllString(str, ""), "");
}

// Count the number of lightning bolt emojis. They take up 2 columns of space so
// that needs to be factored into the length calculation
func countLightningBolts(str string) uint {
    return uint(strings.Count(str, "⚡"));
}

// Find the length of the string (this works for strings with multibyte
// characters)
func strlen(str string) uint {
    // Lightning bolts are reduced down to 1 character using this function, but
    // they actually take up 2 columns of space, so I need to add 1 to our
    // strlen for every lightning bolt.
    return uint(utf8.RuneCountInString(clean(str))) + countLightningBolts(str);
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

// Take the list of examples, and return them as a list of strings
func getExamples(examples []Example) string {
    var ret string = "";
    var padding uint;

    examplesWithLengths := insertLengths(examples);
    maxLength           := maxLength(examplesWithLengths);

    for _, example := range examplesWithLengths {
        padding  = maxLength - example.length + 2;
        ret     += "    " + example.prompt + strings.Repeat(" ", int(padding)) + "# " + example.description;
        ret     += "\n";
    }

    return ret;
}

type Example struct {
    prompt      string;
    description string;
    length      uint;
}


// This function prints out the help for this tool
func help() string {
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
            description: "Auto fetching is in progress. You shouldn't do a pull or fetch when this is happening.",
            prompt: showPrompt(newGitData(GitData{
                fetching: true,
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
                localAhead:   3,
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
                remoteAhead:  3,
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
                remoteAhead:  3,
                localBehind:  3,
                localAhead:   5,
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
        {
            description: "Checked out a tag called my-tag",
            prompt: showPrompt(newGitData(GitData{
                branches: Branches{
                    local: "detached@my-tag",
                },
            })),
        },
    };

    var ret string = "\n";
    ret += "  git-radar - a heads up display for git\n";
    ret += "\n";
    ret += "  examples:\n";
    ret += getExamples(examples);
    ret += "\n";
    ret += "  usage:\n";
    ret += "    git-radar [help|fetch]\n";
    ret += "\n";
    ret += "    fetch  # Fetches your repo asynchronously in the background every 5 mins\n";
    ret += "    help   # Output this help text.\n";
    ret += "\n";
    ret += "  Bash example:\n";
    ret += "    export PS1=\"\\$(git-radar fetch) \"\n";
    ret += "\n";
    ret += "    This will show your current directory and the full git-radar.\n";
    ret += "    As an added benefit, if you are in a repo, it will asynchronously\n";
    ret += "    run `git fetch` every 5 mins, so that you are never out of date.\n";
    ret += "\n";

    return ret;
}
