#!/bin/sh

if [ -f 'src/config.go' ]; then
    echo 'Config file should not be in the repository. Remove it from your commit.';
    exit 1;
fi

if ! cp config.go.example src/config.go; then
    echo "Failed to copy example config in preparation for testing.";
    exit 1;
fi

make clean && make && make test
