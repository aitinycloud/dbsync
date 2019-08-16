GO111MODULE=off

all: build

.PHONY: build
build: dbsync

install:
	go get -v "github.com/go-sql-driver/mysql"
	go get -v "github.com/lib/pq"
	go get -v "cloud.google.com/go"
	go get -v "github.com/mattn/go-oci8"
	go get -v "github.com/patrickmn/go-cache"
	go get -v "github.com/tidwall/gjson"

.PHONY: dbsync
dbsync:
	export GO111MODULE=off ; go build -o dbsync main.go

.PHONY: release
release:
	cp -rf dbsync ./release/
	cp -rf config.json ./release/
	tar -czf dbsync.tar.gz release

.PHONY: test
test:
	go test $(RACE) ./...
