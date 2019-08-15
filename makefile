GO111MODULE=off

all: build

.PHONY: build
build: dbsync

.PHONY: dbsync
dbsync:
	export GO111MODULE=off ; go build -o dbsync ./main.go

.PHONY: release
release:
	cp -rf dbsync ./release/
	cp -rf config.json ./release/
	tar -czf dbsync.tar.gz release

.PHONY: test
test:
	go test $(RACE) ./...