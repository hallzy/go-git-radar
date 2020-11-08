package main

import (
    "fmt"
)

func main() {
    args := getArgs();

    if (len(args) == 0) {
        help();
        return;
    }

    // if not a repo, then nothing to do. Exit silently
    if (!isRepo()) {
        return;
    }

    var isFetch bool = false;
    var shell string = "";

    for _, command := range args {
        switch command {
            case "fetch":
                isFetch = true;
            case "zsh":
                fallthrough;
            case "bash":
                fallthrough;
            case "fish":
                shell = command;
            default:
                panic("Error: [" + command + "] is an unknown option to git-radar");
        }
    }

    if (shell == "") {
        panic("Error: No shell type provided to git-radar");
    }

    if (isFetch) {
        fetch();
    }

    fmt.Print(showPrompt(shell) + "\n");
}
