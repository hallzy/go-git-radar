build: src src/config.go
	$(MAKE) -C src build

test: scratch src tests
	$(MAKE) -C scratch test

scratch:
	mkdir -p scratch && \
		printf "test:\n\tgo test -v -cover" > scratch/Makefile && \
		find . -name "*.go" -exec ln -s ../{} scratch \;

src/config.go:
	$(error "You are missing your comfig file. Please copy config.go.example to src/config.go and make any config changes you want.")
