build: src src/config.go
	$(MAKE) -C src build

test: src tests symlinks
	$(MAKE) -C .scratch test

symlinks:
	find src tests -name "*.go" -exec ln -sf ../{} .scratch \;

src/config.go:
	$(error "You are missing your comfig file. Please copy config.go.example to src/config.go and make any config changes you want.")
