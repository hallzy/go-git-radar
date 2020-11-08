package main

import (
    "fmt"
)

func showRemoteDiverged(behindBy int, aheadBy int) string {
    return fmt.Sprintf(REMOTE_DIVERGED, behindBy, aheadBy);
}

func showRemoteAhead(aheadBy int) string {
    return fmt.Sprintf(REMOTE_AHEAD, aheadBy);
}

func showRemoteBehind(behindBy int) string {
    return fmt.Sprintf(REMOTE_BEHIND, behindBy);
}

func showLocalDiverged(behindBy int, aheadBy int) string {
    return fmt.Sprintf(LOCAL_DIVERGED, behindBy, aheadBy);
}

func showLocalAhead(aheadBy int) string {
    return fmt.Sprintf(LOCAL_AHEAD, aheadBy);
}

func showLocalBehind(behindBy int) string {
    return fmt.Sprintf(LOCAL_BEHIND, behindBy);
}

func showStash(count int) string {
    return fmt.Sprintf(STASH_FORMAT, count);
}

func showUntracked(count int) string {
    if (count == 0) {
        return "";
    }
    return " " + fmt.Sprintf(CHANGES_UNTRACKED, count, "A");
}

func showConflictedUs(count int) string {
    return fmt.Sprintf(CHANGES_CONFLICTED, count, "U");
}

func showConflictedThem(count int) string {
    return fmt.Sprintf(CHANGES_CONFLICTED, count, "T");
}

func showConflictedBoth(count int) string {
    return fmt.Sprintf(CHANGES_CONFLICTED, count, "B");
}

func showUnstagedDeleted(count int) string {
    return fmt.Sprintf(CHANGES_UNSTAGED, count, "D");
}

func showUnstagedModified(count int) string {
    return fmt.Sprintf(CHANGES_UNSTAGED, count, "M");
}

func showUnstagedTypeChanged(count int) string {
    return fmt.Sprintf(CHANGES_UNSTAGED, count, "TC");
}

func showStagedAdded(count int) string {
    return fmt.Sprintf(CHANGES_STAGED, count, "A");
}

func showStagedDeleted(count int) string {
    return fmt.Sprintf(CHANGES_STAGED, count, "D");
}

func showStagedModified(count int) string {
    return fmt.Sprintf(CHANGES_STAGED, count, "M");
}

func showStagedRenamed(count int) string {
    return fmt.Sprintf(CHANGES_STAGED, count, "R");
}

func showStagedCopied(count int) string {
    return fmt.Sprintf(CHANGES_STAGED, count, "C");
}

func showStagedTypeChanged(count int) string {
    return fmt.Sprintf(CHANGES_STAGED, count, "TC");
}

func showConflicted(us int, them int, both int) string {
    var conflicted string = "";

    if (us > 0) {
        conflicted += showConflictedUs(us);
    }

    if (them > 0) {
        conflicted += showConflictedThem(them);
    }

    if (both > 0) {
        conflicted += showConflictedBoth(both);
    }

    if (conflicted == "") {
        return "";
    }
    return " " + conflicted;
}

func showStaged(added int, deleted int, modified int, renamed int, copied int, typeChanged int) string {
    var staged string = "";

    if (added > 0) {
        staged += showStagedAdded(added);
    }

    if (deleted > 0) {
        staged += showStagedDeleted(deleted);
    }

    if (modified > 0) {
        staged += showStagedModified(modified);
    }

    if (renamed > 0) {
        staged += showStagedRenamed(renamed);
    }

    if (copied > 0) {
        staged += showStagedCopied(copied);
    }

    if (typeChanged > 0) {
        staged += showStagedTypeChanged(typeChanged);
    }

    if (staged == "") {
        return "";
    }

    return " " + staged;
}

func showUnstaged(deleted int, modified int, typeChanged int) string {
    var unstaged string = "";

    if (deleted > 0) {
        unstaged += showUnstagedDeleted(deleted);
    }

    if (modified > 0) {
        unstaged += showUnstagedModified(modified);
    }

    if (typeChanged > 0) {
        unstaged += showUnstagedTypeChanged(typeChanged);
    }

    if (unstaged == "") {
        return "";
    }
    return " " + unstaged;
}

func showPrompt(shell string) string {
    remote := getRemoteInfo();
    branch := getBranchInfo();
    local  := getLocalInfo();
    stash  := getStashInfo();
    change := getChangeInfo();

    return fmt.Sprintf(PROMPT_FORMAT, PREFIX, remote, branch, local, stash, change, SUFFIX);
}
