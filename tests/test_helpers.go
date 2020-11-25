package main

// Helper function to test for panics
// test() should throw a panic
// err() is something that happens if test() doesn't panic.
func panicHelper(test func(), err func()) {
    // Recover panics, I don't care about what the panic is though.
    defer func() {
        recover();
    }();

    // This test should throw a panic, which would mean that we recover the
    // panic and don't run the err() function.
    test();

    // test() did not panic, so we are running the error() function
    err();
}
