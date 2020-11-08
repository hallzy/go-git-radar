package main

import (
    "fmt"
)

func gitStr(middle string) string {
    var prefix string = colourString("git:(", "black", true);
    var suffix string = colourString(")", "black", true);

    return prefix + middle + suffix;
}

func masterStr() string {
    return colourString("master", "gray", false);
}

func branchStr(branch string) string {
    return colourString(branch, "gray", false);
}

func untrackedStr(num string) string {
    return num + colourString("A", "gray", true);
}

func addedStagedStr(num string) string {
    return num + colourString("A", "green", true);
}

func modifiedUnstagedStr(num string) string {
    return num + colourString("M", "red", true);
}

func localUpStr(num string) string {
    return num + colourString("↑", "green", true);
}

func fancyMString() string {
    return "\xF0\x9D\x98\xAE";
}

func XFromMasterStr(x string) string {
    return fancyMString() + " " + x + " " + colourString("→", "red", true) + " ";
}

func divergedFromMasterStr(x string, y string) string {
    return fancyMString() + " " + x + " " + colourString("⇄", "yellow", true) + " " + y + " ";
}

func notUpstreamStr() string {
    return "upstream " + colourString("⚡", "red", true) + " ";
}

func detachedStr(hash string) string {
    return colourString("detached@" + hash, "gray", false);
}

func conflictedUsStr(num string) string {
    return num + colourString("U", "yellow", true);
}

func conflictedThemStr(num string) string {
    return num + colourString("T", "yellow", true);
}

func aheadMasterStr(x string) string {
    return fancyMString() + " " + colourString("←", "green", true) + " " + x;
}

func localDivergedStr(x string, y string) string {
    return x + colourString("⇵", "yellow", true) + y;
}

func stashStr(x string) string {
    return x + colourString("≡", "yellow", true);
}

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
