test:
	@echo ''
	@go test -v -coverprofile=coverage.tmp | \
		sed ''/PASS/s//`printf "\033[1;32mPASS\033[0m"`/'' | \
		sed ''/FAIL/s//`printf "\033[1;31mFAIL\033[0m"`/''  | \
		sed ''/SKIP/s//`printf "\033[1;35mSKIP\033[0m"`/''
	@rm testfile-shouldnt-exist.txt || true
	@echo ''

test-report: test
	go tool cover -html=coverage.tmp
