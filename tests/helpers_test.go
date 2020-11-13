package main

import(
    "testing"
    "reflect"
)

func TestParseGitStatus(T *testing.T) {
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
