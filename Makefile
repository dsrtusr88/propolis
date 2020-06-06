GO = GO111MODULE=on go
VERSION=`git describe --tags`

all: fmt check test-coverage build

prepare:
	${GO} get -u github.com/divan/depscheck
	${GO} get github.com/warmans/golocc
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(go env GOPATH)/bin v1.27.0

deps:
	${GO} mod download

fmt:
	${GO} fmt ./...

check: fmt
	golangci-lint run

info: fmt
	depscheck -totalonly -tests .
	golocc .

test-coverage:
	${GO} test -race -coverprofile=coverage.txt -covermode=atomic ./...

clean:
	rm -f propolis
	rm -f propolis-bot
	rm -f propolis_x86
	rm -f propolis_darwin
	rm -f propolis_windows.exe
	rm -f coverage.txt

build:
	cd cmd/propolis;${GO} build -trimpath -ldflags "-X main.Version=${VERSION}" -o ../../propolis;cd ../..
	cd cmd/propolis;GOARCH=386 ${GO} build -ldflags "-X main.Version=${VERSION}" -o ../../propolis_x86;cd ../..
	cd cmd/propolis;GOOS=darwin GOARCH=amd64 ${GO} build -ldflags "-X main.Version=${VERSION}" -o ../../propolis_darwin;cd ../..
	cd cmd/propolis;GOOS=windows GOARCH=amd64 ${GO} build -ldflags "-X main.Version=${VERSION}" -o ../../propolis_windows.exe;cd ../..

build-bot:
	cd cmd/propolis-bot;${GO} build -trimpath -ldflags "-X main.Version=${VERSION}" -o ../../propolis-bot;cd ../..

install:
	cd cmd/propolis;${GO} install -ldflags "-X main.Version=${VERSION}";cd ../..





