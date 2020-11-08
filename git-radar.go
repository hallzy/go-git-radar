package main

import (
    "fmt"
)

func main() {
    Git := getGitData();

    // if not a repo, then nothing to do. Exit silently
    if (!Git.isRepo) {
        return;
    }

    args := getArgs();
    var isFetch bool = false;

    for _, command := range args {
        switch command {
            case "fetch":
                isFetch = true;
            case "help":
                help();
                return;
            default:
                panic("Error: [" + command + "] is an unknown option to git-radar");
        }
    }

    if (isFetch) {
        fetch(Git.dotGit);
    }

    fmt.Print(showPrompt(Git) + "\n");
}
