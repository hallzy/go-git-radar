package main

import(
    "testing"
    "reflect"
)

// Test clean(){{{
func TestClean(T *testing.T) {
    inputExpected := map[string]string {
        "This is a test":                                  "This is a test",
        "\x01\033[1;31m\x02This\x01\033[0m\x02 is a test": "This is a test",
    }

    for input, expected := range inputExpected {
        output := clean(input);
        if (output != expected) {
            T.Errorf("clean(): Got [%s], expected [%s] for input [%s]", output, expected, input);
        }
    }
}
// }}}
// Test strlen(){{{
func TestStrlen(T *testing.T) {
    inputExpected := map[string]uint {
        "This is a test":    14,
        "git:(⚡my-branch)": 16,
    }

    for input, expected := range inputExpected {
        output := strlen(input);
        if (output != expected) {
            T.Errorf("strlen(): Got [%d], expected [%d] for input [%s]", output, expected, input);
        }
    }
}
// }}}
// Test insertLengths(){{{
func TestInsertLengths(T *testing.T) {
    input := []Example{
        Example{ prompt: "this is a test prompt" },
        Example{ prompt: "git:(⚡my-branch)" },
        Example{ prompt: "\x01\033[1;31m\x02This\x01\033[0m\x02 is a test" },
    };
    expected := []Example{
        Example{ prompt: "this is a test prompt",                           length: 21 },
        Example{ prompt: "git:(⚡my-branch)",                               length: 16 },
        Example{ prompt: "\x01\033[1;31m\x02This\x01\033[0m\x02 is a test", length: 14 },
    };

    output := insertLengths(input);
    if (!reflect.DeepEqual(output, expected)) {
        T.Errorf("strlen(): Got [%+v], expected [%+v] for input [%+v]", output, expected, input);
    }
}
// }}}
// Test maxLength(){{{
func TestMaxLength(T *testing.T) {
    input := []Example{
        Example{ prompt: "this is a test prompt",                           length: 21 },
        Example{ prompt: "git:(⚡my-branch)",                               length: 16 },
        Example{ prompt: "\x01\033[1;31m\x02This\x01\033[0m\x02 is a test", length: 14 },
    };

    output := maxLength(input);
    if (output != 21) {
        T.Errorf("strlen(): Got [%d], expected [21] for input [%+v]", output, input);
    }
}
// }}}
// Test help(){{{
func TestHelp(T *testing.T) {
    help();
    // There isn't much to test here. It prints out a fixed string that is
    // always the same. Basically we are just running this to artificially bump
    // up the coverage.
}
// }}}
