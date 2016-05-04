REF := $(shell if git diff --quiet HEAD; then git rev-parse --verify HEAD --short; else echo "scratch"; fi)
BUILD := $(shell date -u +%Y%m%d.%H%M%S)

grafanactl: *.go
	go build -ldflags "-X main.REF=$(REF) -X main.BUILD=$(BUILD)"
