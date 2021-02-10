test:
	@echo ''
	@test -f not-writable.txt || touch not-writable.txt
	@test ! -w not-writable.txt || chmod a-w not-writable.txt
	@rm testfile-shouldnt-exist.txt > /dev/null 2>&1 || true
	@go test -v -coverprofile=coverage.tmp | \
		sed ''/PASS/s//`printf "\033[1;32mPASS\033[0m"`/'' | \
		sed ''/FAIL/s//`printf "\033[1;31mFAIL\033[0m"`/''  | \
		sed ''/SKIP/s//`printf "\033[1;35mSKIP\033[0m"`/''
	@rm testfile-shouldnt-exist.txt > /dev/null 2>&1 || true
	@echo ''

test-report: test
	sed -i "s/_$$(pwd|sed 's@/@\\/@g')/./g" coverage.tmp
	go tool cover -html=coverage.tmp
