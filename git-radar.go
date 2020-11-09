package main

import (
    "fmt"
)

func main() {
    // Get important git information about this current repo we are in
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

    // Print out the prompt with the given Git data
    fmt.Print(showPrompt(Git) + "\n");
}
