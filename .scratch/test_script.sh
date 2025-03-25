#!/bin/sh

echo ''

if [ ! -f 'not-writable.txt' ]; then
    touch not-writable.txt
fi

if [ -w 'not-writable.txt' ]; then
    chmod a-w not-writable.txt
fi

rm testfile-shouldnt-exist.txt > /dev/null 2>&1

results="$(GO111MODULE=off go test -v -coverprofile=coverage.tmp)"
testSucceeded=$?


echo "$results" | \
    sed ''/PASS/s//`printf "\033[1;32mPASS\033[0m"`/'' | \
    sed ''/FAIL/s//`printf "\033[1;31mFAIL\033[0m"`/''  | \
    sed ''/SKIP/s//`printf "\033[1;35mSKIP\033[0m"`/''

rm testfile-shouldnt-exist.txt > /dev/null 2>&1

echo ''

exit $testSucceeded
