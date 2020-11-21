package main

import(
    "testing"
    "reflect"
)

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
            T.Errorf("str2int(): string to uint conversion test failed. Got [%d], expected [%d] for input [%s]", output, expected, input);
        }
    }

    // Run all tests that should fail
    for _, input := range inputPanic {
        func() {
            // Recover panics, I don't care about what the panic is though.
            defer func() {
                recover();
            }();

            // Attempt to convert
            output := str2int(input);

            // If I actually get to this point, then the function execution
            // succeeded because the defer function didn't run. This is an error
            // because all tests that are run here should panic
            T.Errorf("str2int(): string to uint conversion test failed. input [%s] should have panicked, but got [%d]", input, output);
        }();
    }
}

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
            T.Errorf("int2str(): uint to string conversion test failed. Got [%s], expected [%s] for input [%d]", output, expected, input);
        }
    }
}

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
            T.Errorf("trim(): string trim test failed. Got [%s], expected [%s] for input [%s]", output, expected, input);
        }
    }
}

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
            T.Errorf("parseGitStatus(): Partially staged file test failed. Got [%+v], expected [%+v]", got, expected);
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
            stagedModified: 2,
            stagedAdded: 1,
            stagedDeleted: 2,
            stagedRenamed: 3,
            stagedCopied: 1,
            stagedTypeChanged: 1,

            unstagedModified: 4,
            unstagedDeleted: 2,
            unstagedTypeChanged: 1,

            conflictUs: 1,
            conflictThem: 1,
            conflictBoth: 5,

            untracked: 3,
        };

        got := parseGitStatus(input);

        if (!reflect.DeepEqual(got, expected)) {
            T.Errorf("parseGitStatus(): Combination of all types. Got [%+v], expected [%+v]", got, expected);
        }
    }
    // }}}
    filesWithStagedAndUnstagedChanges();
    combinationOfAll();
}

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
            T.Errorf("countNewLines(): count newlines test failed. Got [%d], expected [%d] for input [%s]", output, expected, input);
        }
    }
}
