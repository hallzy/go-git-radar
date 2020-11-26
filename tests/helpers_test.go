package main

import(
    "testing"
    "reflect"
)

// Test str2int{{{
func TestStr2int(T *testing.T) {
    inputExpected := map[string]uint {
        "10":         10,
        "0":          0,
        "1":          1,
        "4294967295": 4294967295,
    }

    // These should produce panics
    inputPanic := []string {
        "4294967296",
        "1.2",
        "-1",
        "1.2e4",
    }

    // Run all tests that should pass
    for input, expected := range inputExpected {
        output := str2int(input);
        if (output != expected) {
            T.Errorf("str2int(): Got [%d], expected [%d] for input [%s]", output, expected, input);
        }
    }

    // Run all tests that should fail
    var output uint;
    for _, input := range inputPanic {
        panicHelper(
            func() {
                output = str2int(input);
            },
            func() {
                T.Errorf("str2int(): Input [%s] should have panicked, but got [%d]", input, output);
            },
        );
    }
}
// }}}
// Test int2str{{{
func TestInt2str(T *testing.T) {
    inputExpected := map[uint]string {
        10:         "10",
        0:          "0",
        1:          "1",
        4294967295: "4294967295",
        1.2e4:      "12000",
    }

    for input, expected := range inputExpected {
        output := int2str(input);
        if (output != expected) {
            T.Errorf("int2str(): Got [%s], expected [%s] for input [%d]", output, expected, input);
        }
    }
}
// }}}
// Test insertData(){{{
func TestInsertData(T *testing.T) {
    str := "%%HELLO%% | %%HELLO%% | %%BYE%% | %HELLO% | HELLO";
    data := FormatData{
        "HELLO": "Hello",
        "BYE": "Bye",
    };

    expected := "Hello | Hello | Bye | %HELLO% | HELLO";

    output := insertData(str, data);
    if (output != expected) {
        T.Errorf("int2str(): Got [%s], expected [%s] for input [%s, %+v]", output, expected, str, data);
    }
}
// }}}
// Test trim{{{
func TestTrim(T *testing.T) {
    inputExpected := map[string]string {
        "Hello there":               "Hello there",
        " Hello there":              "Hello there",
        " Hello there ":             "Hello there",
        "    Hello there  ":         "Hello there",
        " \t  Hello there  \n\t\n ": "Hello there",
        " Hello   there ":           "Hello   there",
    }

    for input, expected := range inputExpected {
        output := trim(input);
        if (output != expected) {
            T.Errorf("trim(): Got [%s], expected [%s] for input [%s]", output, expected, input);
        }
    }
}
// }}}
// Test parseGitStatus{{{
func TestParseGitStatus(T *testing.T) {
    // 2 files have staged and unstaged changes, one file has unstaged changes {{{
    filesWithStagedAndUnstagedChanges := func() {
        got := parseGitStatus([]string {
            "MM file1",
            " M file2",
            "MM file3",
        });

        expected := GitStatus{
            stagedModified:      2,
            unstagedModified:    3,
        };

        if (!reflect.DeepEqual(got, expected)) {
            T.Errorf("parseGitStatus(): Got [%+v], expected [%+v]", got, expected);
        }
    }
    // }}}
    // Combination of every status type {{{
    combinationOfAll := func() {
        input := []string {
            "M   staged modified 1",
            " D  unstaged deleted 1",
            "xD  unstaged deleted 2",
            " T  unstaged type change 1",
            "DD  conflict both 4",
            "UU  conflict both 1",
            "A   staged new file 1",
            "R   staged renamed 2",
            "Dx  staged deleted 2",
            " M  unstaged modified 2",
            "D   staged deleted 1",
            "Mx  staged modified 2",
            "C   staged copied 1",
            "??  untracked 3",
            " M  unstaged modified 1",
            "DD  conflict both 5",
            "U   conflict them 1",
            "R   staged renamed 1",
            " U  conflict us 1",
            "UU  conflict both 2",
            "R   staged renamed 3",
            "T   staged type change 1",
            "xM  unstaged modified 3",
            " M  unstaged modified 4",
            "AA  conflict both 3",
            "??  untracked 1",
            "??  untracked 2",
        };

        expected := GitStatus{
            stagedModified:    2,
            stagedAdded:       1,
            stagedDeleted:     2,
            stagedRenamed:     3,
            stagedCopied:      1,
            stagedTypeChanged: 1,

            unstagedModified:    4,
            unstagedDeleted:     2,
            unstagedTypeChanged: 1,

            conflictUs:   1,
            conflictThem: 1,
            conflictBoth: 5,

            untracked: 3,
        };

        got := parseGitStatus(input);

        if (!reflect.DeepEqual(got, expected)) {
            T.Errorf("parseGitStatus(): Got [%+v], expected [%+v]", got, expected);
        }
    }
    // }}}
    filesWithStagedAndUnstagedChanges();
    combinationOfAll();
}
// }}}
// Test countLines{{{
func TestCountLines(T *testing.T) {
    inputExpected := map[string]uint {
        "":                   0,
        "\n\n":               2,
        "Testing\n\nTesting": 2,
        "\r\n\r\n\r\n":       3,
    }

    for input, expected := range inputExpected {
        output := countNewLines(input);
        if (output != expected) {
            T.Errorf("countNewLines(): Got [%d], expected [%d] for input [%s]", output, expected, input);
        }
    }
}
// }}}
// Test getFullRemote{{{
func TestGetFullRemote(T *testing.T) {
    inputExpected := map[RemoteBranch]string {
        RemoteBranch{remote: "origin",     branch: "master"}:             "origin/master",
        RemoteBranch{remote: "origin",     branch: "anything"}:           "origin/anything",
        RemoteBranch{remote: "origin",     branch: "anything/something"}: "origin/anything/something",
        RemoteBranch{remote: "not-origin", branch: "master"}:             "not-origin/master",
        RemoteBranch{remote: "not-origin", branch: "anything"}:           "not-origin/anything",
        RemoteBranch{remote: "not-origin", branch: "anything/something"}: "not-origin/anything/something",
        RemoteBranch{remote: "",           branch: ""}:                   "",
        RemoteBranch{remote: "origin",     branch: ""}:                   "",
        RemoteBranch{remote: "not-origin", branch: ""}:                   "",
        RemoteBranch{remote: "",           branch: "somebranch"}:         "",
    }

    for input, expected := range inputExpected {
        output := getFullRemote(input);
        if (output != expected) {
            T.Errorf("getFullRemote(): Got [%s], expected [%s] for input [%+v]", output, expected, input);
        }
    }
}
// }}}
// Test getFullName{{{
func TestGetFullName(T *testing.T) {
    inputExpected := map[RemoteBranch]string {
        RemoteBranch{remote: "origin",     branch: "master"}:             "master",
        RemoteBranch{remote: "origin",     branch: "anything"}:           "anything",
        RemoteBranch{remote: "origin",     branch: "anything/something"}: "anything/something",
        RemoteBranch{remote: "not-origin", branch: "master"}:             "not-origin/master",
        RemoteBranch{remote: "not-origin", branch: "anything"}:           "not-origin/anything",
        RemoteBranch{remote: "not-origin", branch: "anything/something"}: "not-origin/anything/something",
        RemoteBranch{remote: "",           branch: ""}:                   "",
        RemoteBranch{remote: "origin",     branch: ""}:                   "",
        RemoteBranch{remote: "not-origin", branch: ""}:                   "",
        RemoteBranch{remote: "",           branch: "somebranch"}:         "",
    }

    for input, expected := range inputExpected {
        output := getFullName(input);
        if (output != expected) {
            T.Errorf("getFullName(): Got [%s], expected [%s] for input [%+v]", output, expected, input);
        }
    }
}
// }}}
// Test getRemoteInfo{{{
func TestGetRemoteInfo(T *testing.T) {
    type TestRemoteInfoType struct {
        branches     Branches;
        remoteBehind uint;
        remoteAhead  uint;
    }

    inputExpected := map[TestRemoteInfoType]string {
        // Test the remote same
        TestRemoteInfoType{
            remoteBehind: 0, remoteAhead: 0, branches: Branches {
                local: "",
                remote: RemoteBranch{remote: "origin", branch: "master"},
                parent: RemoteBranch{remote: "origin", branch: "master"},
            },
        } : insertData(REMOTE_SAME, FormatData{
            "REMOTE_BRANCH": "master",
        }),

        // Test remote equal
        TestRemoteInfoType{
            remoteBehind: 0, remoteAhead: 0, branches: Branches {
                local: "",
                remote: RemoteBranch{remote: "origin", branch: "branch"},
                parent: RemoteBranch{remote: "origin", branch: "master"},
            },
        } : insertData(REMOTE_EQUAL, FormatData{
            "PARENT_REMOTE_BRANCH": "master",
            "REMOTE_BRANCH":        "branch",
        }),

        // test remote behind
        TestRemoteInfoType{
            remoteBehind: 1, remoteAhead: 0, branches: Branches {
                local: "",
                remote: RemoteBranch{remote: "origin",   branch: "branch"},
                parent: RemoteBranch{remote: "upstream", branch: "master"},
            },
        } : insertData(REMOTE_BEHIND, FormatData{
            "PARENT_REMOTE_BRANCH": "upstream/master",
            "REMOTE_BEHIND":        "1",
            "REMOTE_BRANCH":        "branch",
        }),

        // test remote ahead
        TestRemoteInfoType{
            remoteBehind: 0, remoteAhead: 3, branches: Branches {
                local: "",
                remote: RemoteBranch{remote: "origin",   branch: "branch"},
                parent: RemoteBranch{remote: "upstream", branch: "master"},
            },
        } : insertData(REMOTE_AHEAD, FormatData{
            "PARENT_REMOTE_BRANCH": "upstream/master",
            "REMOTE_AHEAD":        "3",
            "REMOTE_BRANCH":        "branch",
        }),

        // test remote diverged
        TestRemoteInfoType{
            remoteBehind: 1, remoteAhead: 2, branches: Branches {
                local: "",
                remote: RemoteBranch{remote: "origin",   branch: "branch"},
                parent: RemoteBranch{remote: "upstream", branch: "master"},
            },
        } : insertData(REMOTE_DIVERGED, FormatData{
            "PARENT_REMOTE_BRANCH": "upstream/master",
            "REMOTE_BEHIND":        "1",
            "REMOTE_AHEAD":         "2",
            "REMOTE_BRANCH":        "branch",
        }),
    }

    for input, expected := range inputExpected {
        output := getRemoteInfo(input.branches, input.remoteBehind, input.remoteAhead);
        if (output != expected) {
            T.Errorf("getRemoteInfo(): Got [%s], expected [%s] for input [%+v]", output, expected, input);
        }
    }
}
// }}}
// Test getLocalInfo{{{
func TestGetLocalInfo(T *testing.T) {
    type TestDiff struct {
        ahead  uint;
        behind uint;
    }
    inputExpected := map[TestDiff]string {
        TestDiff{ahead: 0, behind: 0}: "",
        TestDiff{ahead: 1, behind: 2}: insertData(LOCAL_DIVERGED, FormatData{
            "LOCAL_BEHIND":        "2",
            "LOCAL_AHEAD":         "1",
        }),
        TestDiff{ahead: 0, behind: 3}: insertData(LOCAL_BEHIND, FormatData{ "LOCAL_BEHIND": "3", }),
        TestDiff{ahead: 2, behind: 0}: insertData(LOCAL_AHEAD, FormatData{ "LOCAL_AHEAD": "2", }),
    }

    for input, expected := range inputExpected {
        output := getLocalInfo(input.behind, input.ahead);
        if (output != expected) {
            T.Errorf("getLocalInfo(): Got [%s], expected [%s] for input [%+v]", output, expected, input);
        }
    }
}
// }}}
// Test showUntracked{{{
func TestShowUntracked(T *testing.T) {
    inputExpected := map[GitStatus]string {
        GitStatus{
            untracked: 0,

            stagedAdded:         5,
            stagedDeleted:       5,
            stagedModified:      5,
            stagedRenamed:       5,
            stagedCopied:        5,
            stagedTypeChanged:   5,
            unstagedDeleted:     5,
            unstagedModified:    5,
            unstagedTypeChanged: 5,
            conflictUs:          5,
            conflictThem:        5,
            conflictBoth:        5,
        }: "",
        GitStatus{
            untracked: 2,

            stagedAdded:         5,
            stagedDeleted:       5,
            stagedModified:      5,
            stagedRenamed:       5,
            stagedCopied:        5,
            stagedTypeChanged:   5,
            unstagedDeleted:     5,
            unstagedModified:    5,
            unstagedTypeChanged: 5,
            conflictUs:          5,
            conflictThem:        5,
            conflictBoth:        5,
        }: " " + insertData(CHANGES_UNTRACKED, FormatData{ "COUNT": "2" }),
    };

    for input, expected := range inputExpected {
        output := showUntracked(input);
        if (output != expected) {
            T.Errorf("showUntracked(): Got [%s], expected [%s] for input [%+v]", output, expected, input);
        }
    }
}
// }}}
// Test showConflicted{{{
func TestShowConflicted(T *testing.T) {
    inputExpected := map[GitStatus]string {
        GitStatus{
            conflictUs:   0,
            conflictThem: 0,
            conflictBoth: 0,

            untracked:           5,
            stagedAdded:         5,
            stagedDeleted:       5,
            stagedModified:      5,
            stagedRenamed:       5,
            stagedCopied:        5,
            stagedTypeChanged:   5,
            unstagedDeleted:     5,
            unstagedModified:    5,
            unstagedTypeChanged: 5,
        }: "",
        GitStatus{
            conflictUs:   1,
            conflictThem: 0,
            conflictBoth: 0,

            untracked:           5,
            stagedAdded:         5,
            stagedDeleted:       5,
            stagedModified:      5,
            stagedRenamed:       5,
            stagedCopied:        5,
            stagedTypeChanged:   5,
            unstagedDeleted:     5,
            unstagedModified:    5,
            unstagedTypeChanged: 5,
        }: " " + insertData(CHANGES_CONFLICTED, FormatData{ "COUNT": "1", "SYMBOL": CONFLICT_US_SYM }),
        GitStatus{
            conflictUs:   0,
            conflictThem: 2,
            conflictBoth: 3,

            untracked:           5,
            stagedAdded:         5,
            stagedDeleted:       5,
            stagedModified:      5,
            stagedRenamed:       5,
            stagedCopied:        5,
            stagedTypeChanged:   5,
            unstagedDeleted:     5,
            unstagedModified:    5,
            unstagedTypeChanged: 5,
        }: " " + insertData(CHANGES_CONFLICTED, FormatData{ "COUNT": "2", "SYMBOL": CONFLICT_THEM_SYM }) + insertData(CHANGES_CONFLICTED, FormatData{ "COUNT": "3", "SYMBOL": CONFLICT_BOTH_SYM }),
    }

    for input, expected := range inputExpected {
        output := showConflicted(input);
        if (output != expected) {
            T.Errorf("showConflicted(): Got [%s], expected [%s] for input [%+v]", output, expected, input);
        }
    }
}
// }}}
// Test showStaged{{{
func TestShowStaged(T *testing.T) {
    inputExpected := map[GitStatus]string {
        GitStatus{
            stagedAdded:       0,
            stagedDeleted:     0,
            stagedModified:    0,
            stagedRenamed:     0,
            stagedCopied:      0,
            stagedTypeChanged: 0,

            conflictBoth:        5,
            conflictThem:        5,
            conflictUs:          5,
            unstagedDeleted:     5,
            unstagedModified:    5,
            unstagedTypeChanged: 5,
            untracked:           5,
        }: "",
        GitStatus{
            stagedAdded:       1,
            stagedDeleted:     2,
            stagedModified:    0,
            stagedRenamed:     0,
            stagedCopied:      0,
            stagedTypeChanged: 0,

            conflictBoth:        5,
            conflictThem:        5,
            conflictUs:          5,
            unstagedDeleted:     5,
            unstagedModified:    5,
            unstagedTypeChanged: 5,
            untracked:           5,
        }: " " + insertData(CHANGES_STAGED, FormatData{ "COUNT": "1", "SYMBOL": STAGED_ADDED_SYM }) + insertData(CHANGES_STAGED, FormatData{ "COUNT": "2", "SYMBOL": STAGED_DELETED_SYM }),
        GitStatus{
            stagedAdded:       0,
            stagedDeleted:     0,
            stagedModified:    3,
            stagedRenamed:     4,
            stagedCopied:      0,
            stagedTypeChanged: 0,

            conflictBoth:        5,
            conflictThem:        5,
            conflictUs:          5,
            unstagedDeleted:     5,
            unstagedModified:    5,
            unstagedTypeChanged: 5,
            untracked:           5,
        }: " " + insertData(CHANGES_STAGED, FormatData{ "COUNT": "3", "SYMBOL": STAGED_MODIFIED_SYM }) + insertData(CHANGES_STAGED, FormatData{ "COUNT": "4", "SYMBOL": STAGED_RENAMED_SYM }),
        GitStatus{
            stagedAdded:       0,
            stagedDeleted:     0,
            stagedModified:    0,
            stagedRenamed:     0,
            stagedCopied:      6,
            stagedTypeChanged: 7,

            conflictBoth:        5,
            conflictThem:        5,
            conflictUs:          5,
            unstagedDeleted:     5,
            unstagedModified:    5,
            unstagedTypeChanged: 5,
            untracked:           5,
        }: " " + insertData(CHANGES_STAGED, FormatData{ "COUNT": "6", "SYMBOL": STAGED_COPIED_SYM }) + insertData(CHANGES_STAGED, FormatData{ "COUNT": "7", "SYMBOL": STAGED_TYPE_CHANGED_SYM }),
    }

    for input, expected := range inputExpected {
        output := showStaged(input);
        if (output != expected) {
            T.Errorf("showStaged(): Got [%s], expected [%s] for input [%+v]", output, expected, input);
        }
    }
}
// }}}
// Test showUnstaged{{{
func TestShowUnstaged(T *testing.T) {
    inputExpected := map[GitStatus]string {
        GitStatus{
            unstagedDeleted:     0,
            unstagedModified:    0,
            unstagedTypeChanged: 0,
            untracked:           0,

            conflictBoth:      5,
            conflictThem:      5,
            conflictUs:        5,
            stagedAdded:       5,
            stagedCopied:      5,
            stagedDeleted:     5,
            stagedModified:    5,
            stagedRenamed:     5,
            stagedTypeChanged: 5,
        }: "",
        GitStatus{
            unstagedDeleted:     1,
            unstagedModified:    2,
            unstagedTypeChanged: 0,
            untracked:           0,

            conflictBoth:      5,
            conflictThem:      5,
            conflictUs:        5,
            stagedAdded:       5,
            stagedCopied:      5,
            stagedDeleted:     5,
            stagedModified:    5,
            stagedRenamed:     5,
            stagedTypeChanged: 5,
        }: " " + insertData(CHANGES_UNSTAGED, FormatData{ "COUNT": "1", "SYMBOL": UNSTAGED_DELETED_SYM }) + insertData(CHANGES_UNSTAGED, FormatData{ "COUNT": "2", "SYMBOL": UNSTAGED_MODIFIED_SYM }),
        GitStatus{
            unstagedDeleted:     4,
            unstagedModified:    0,
            unstagedTypeChanged: 3,

            conflictBoth:      5,
            conflictThem:      5,
            conflictUs:        5,
            stagedAdded:       5,
            stagedCopied:      5,
            stagedDeleted:     5,
            stagedModified:    5,
            stagedRenamed:     5,
            stagedTypeChanged: 5,
            untracked:         5,
        }: " " + insertData(CHANGES_UNSTAGED, FormatData{ "COUNT": "4", "SYMBOL": UNSTAGED_DELETED_SYM }) + insertData(CHANGES_UNSTAGED, FormatData{ "COUNT": "3", "SYMBOL": UNSTAGED_TYPE_CHANGED_SYM }),
    }

    for input, expected := range inputExpected {
        output := showUnstaged(input);
        if (output != expected) {
            T.Errorf("showUnstaged(): Got [%s], expected [%s] for input [%+v]", output, expected, input);
        }
    }
}
// }}}
// Test newGitData{{{
func TestNewGitData(T *testing.T) {
    inputExpected := map[GitData]GitData {
        GitData{ }: GitData{
            branches: Branches{
                local: "<unset>",
                parent: RemoteBranch{ remote: "origin", branch: "master", },
            },
        },
    }

    for input, expected := range inputExpected {
        output := newGitData(input);
        if (output != expected) {
            T.Errorf("newGitData(): Got [%+v], expected [%+v] for input [%+v]", output, expected, input);
        }
    }
}
// }}}
// Test isRepo(){{{
func TestIsRepo(T *testing.T) {
    type IsRepoData struct {
        cwd    string;
        dotGit string;
    }

    inputExpected := map[IsRepoData]bool {
        // If . is the .git folder ,then that means we are in the .git folder,
        // so the cwd can be ignored and this is not a repo
        IsRepoData{cwd: "/home/test/repo/hello",               dotGit: "."}:                          false,

        // If no git folder found, then not a git repo
        IsRepoData{cwd: "/home/test/repo/hello",               dotGit: ""}:                           false,

        // .git folder is in a subdirectory, which shouldn't happen
        IsRepoData{cwd: "/home/test/repo",                     dotGit: "/home/test/repo/hello/.git"}: false,

        // We are inside the .git folder, which should be ignored
        IsRepoData{cwd: "/home/test/repo/hello/.git/branches", dotGit: "/home/test/repo/hello/.git"}: false,

        // We are in a repo with the .git folder in the cwd which means we are
        // in the root of the repo.
        IsRepoData{cwd: "/home/test/repo/hello",               dotGit: ".git"}:                       true,

        // We are in the root of the repo, but a full .git path was given
        IsRepoData{cwd: "/home/test/repo/hello",               dotGit: "/home/test/repo/hello/.git"}: true,

        // We are in a subdirectory of the repo
        IsRepoData{cwd: "/home/test/repo/hello/somedir",       dotGit: "/home/test/repo/hello/.git"}: true,
    }

    for input, expected := range inputExpected {
        output := isRepo(input.cwd, input.dotGit);
        if (output != expected) {
            T.Errorf("newGitData(): Got [%t], expected [%t] for input [%s, %s]", output, expected, input.cwd, input.dotGit);
        }
    }
}
// }}}
