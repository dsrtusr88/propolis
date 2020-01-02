GO = GO111MODULE=on go

all: fmt check test-coverage build

prepare:
	${GO} get -u github.com/divan/depscheck
	${GO} get github.com/warmans/golocc
	${GO} install github.com/golangci/golangci-lint/cmd/golangci-lint

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
	rm -f propolis_darwin
	rm -f propolis_windows.exe
	rm -f coverage.txt

build:
	${GO} build -v -o propolis
	GOOS=darwin GOARCH=amd64 ${GO} build -v -o propolis_darwin
	GOOS=windows GOARCH=amd64 ${GO} build -v -o propolis_windows.exe

install:
	${GO} install -v ./...





