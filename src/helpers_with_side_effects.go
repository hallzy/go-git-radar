package main

// =============================================================================
// ALL FUNCTIONS IN THIS FILE HAVE SIDE EFFECTS AND/OR ARE DEPENDENT ON SOME
// EXTERNAL STATE
// =============================================================================

import (
    "os"
    "os/exec"
    "io/ioutil"
    "time"
)

// Get command line arguments
func getArgs() []string {
    return os.Args[1:];
}

// Get unix timestamp
func now() uint {
    return uint(time.Now().Unix());
}

func fileExists(file string) bool {
    _, err := os.Stat(file);

    return ! os.IsNotExist(err);
}

func fileRead(file string) string {
    ret, err := ioutil.ReadFile(file);

    if (err != nil) {
        panic(err.Error());
    }

    return string(ret);
}

func fileWrite(file string, data string) bool {
    err := ioutil.WriteFile(file, []byte(data), 0664);

    if (err != nil) {
        panic(err.Error());
    }

    return true;
}

func getCwd() string {
    cwd, _ := runCmd("pwd");
    return cwd;
}

// Run a command, but don't wait for it to finish
func runCmdConcurrent(cmdStr string) {
    cmd := exec.Command("/bin/sh", "-c", cmdStr);
    err := cmd.Start();

    if (err != nil) {
        panic("Concurrent command [" + cmdStr + "] failed.");
    }
}

// Run a command, but don't trim the output
func runCmdNoTrim(cmdStr string) (string, error) {
    cmd := exec.Command("/bin/sh", "-c", cmdStr);
    out, err := cmd.Output();

    return string(out), err;
}

// Run a command, trim the output
func runCmd(cmdStr string) (string, error) {
    o, e := runCmdNoTrim(cmdStr);

    return trim(o), e;
}
