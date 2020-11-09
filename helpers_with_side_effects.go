package main

// =============================================================================
// ALL FUNCTIONS IN THIS FILE HAVE SIDE EFFECTS AND/OR ARE DEPENDENT ON SOME
// EXTERNAL STATE
// =============================================================================

import (
    "os"
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

