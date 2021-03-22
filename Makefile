build: src src/config.go
	$(MAKE) -C src build

test: .scratch src tests symlinks
	$(MAKE) -C .scratch test

test-report: test
	$(MAKE) -C .scratch test-report

symlinks:
	find src tests -name "*.go" -exec ln -sf ../{} .scratch \;

src/config.go:
	$(error "You are missing your comfig file. Please copy config.go.example to src/config.go and make any config changes you want.")

.scratch:
	$(error ".scratch folder is missing. It must be downloaded from GitHub or recreated in some way before tests can run.")

src:
	$(error "Source files are missing. Redownload from GitHub")

tests:
	$(error "Test files are missing. Redownload from GitHub")

clean:
	rm -f git-radar
