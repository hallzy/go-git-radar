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
    // A case where 2 files have both staged and unstaged changes, while one
    // file just has unstaged changes
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

    filesWithStagedAndUnstagedChanges();
}
