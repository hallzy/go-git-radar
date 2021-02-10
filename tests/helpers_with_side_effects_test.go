package main

import(
    "testing"
    "strings"
)

var skipCases = map[string]string{};

// Test getArgs(){{{
func TestGetArgs(T *testing.T) {
    // The output is non deterministic, but it should have 3 elements
    output := getArgs();

    length := len(output);
    if (length != 2) {
        T.Errorf("getArgs(): Expected 2 element, got %d", length);
    }
}
// }}}
// Test now(){{{
func TestNow(T *testing.T) {
    // This is roughly Nov 24th 2020 at 2:50pm PST, which is the time that I am
    // currently writing this function. So anything after now should be higher
    // than that.
    var fixedTime uint = 1606258204;

    output := now();
    if (fixedTime >= output) {
        T.Errorf("now(): Got [%d], expected to be higher than [%d]", output, fixedTime);
    }
}
// }}}
// Test fileExists(){{{
func TestFileExists(T *testing.T) {
    inputExpected := map[string]bool {
        "git-radar.go":                         true,
        "Makefile":                             true,
        "/a/file/that/definitely/doesnt/exist": false,
    }

    for input, expected := range inputExpected {
        output := fileExists(input);
        if (output != expected) {
            skipCases["TestFileRead"]  = "TestFileRead(): Skipped because TestFileExists() failed.";
            skipCases["TestFileWrite"] = "TestFileWrite(): Skipped because TestFileExists() failed.";

            T.Errorf("clean(): Got [%t], expected [%t] for input [%s]", output, expected, input);
        }
    }
}
// }}}
// Test fileRead(){{{
func TestFileRead(T *testing.T) {
    if val, ok := skipCases["TestFileRead"]; ok {
        T.Skip(val);
    }

    output := fileRead("Makefile");
    if (len(output) <= 100) {
        skipCases["TestFileWrite"] = "TestFileWrite(): Skipped because TestFileRead() failed.";

        T.Errorf("fileRead(): Failed");
    }

    panicHelper(
        func() {
            fileRead("/a/file/that/definitely/doesnt/exist");
        },
        func() {
            T.Errorf("fileRead(): Panic test failed.");
        },
    );
}
// }}}
// Test fileWrite(){{{
func TestFileWrite(T *testing.T) {
    if val, ok := skipCases["TestFileWrite"]; ok {
        T.Skip(val);
    }

    file := "testfile-shouldnt-exist.txt";

    if (fileExists(file)) {
        T.Errorf("fileWrite(): test file [.scratch/" + file + "] already exists. Please. It should be automatically removed by the Makefile");
        return;
    }

    toWrite := "This is a test file\n" + int2str(now());

    ret := fileWrite(file, toWrite);
    if (ret == false) {
        T.Errorf("fileWrite(): Failed to write file");
        return;
    }

    contents := fileRead(file);
    if (contents != toWrite) {
        T.Errorf("fileWrite(): Text read from file is not what was written");
        return;
    }

    panicHelper(
        func() {
            // This file doesn't exist, and creating it would require modify
            fileWrite("/a/file/that/definitely/doesnt/exist", "data");
        },
        func() {
            T.Errorf("fileRead(): Panic test failed.");
        },
    );
}
// }}}
// Test getCwd(){{{
func TestGetCwd(T *testing.T) {
    cwd := getCwd();

    // Can't do much more than check that the end of the cwd is .scratch
    if (!strings.HasSuffix(cwd, ".scratch")) {
        T.Errorf("getCwd(): Failed to get cwd");
        return;
    }
}
// }}}
