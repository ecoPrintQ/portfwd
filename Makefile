EXECUTABLE=portfwd
WINDOWS=$(EXECUTABLE)_windows_amd64.exe
LINUX=$(EXECUTABLE)_linux_amd64
DARWIN=$(EXECUTABLE)_darwin_amd64
VERSION=$(shell git describe --tags)
DATE=$(shell date +%FT%T%z)

build: windows linux darwin ## Build binaries
	go mod download
	@echo version: $(VERSION)

windows: $(WINDOWS) ## Build for Windows

linux: $(LINUX) ## Build for Linux

darwin: $(DARWIN) ## Build for Darwin (macOS)

$(WINDOWS):
	env GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-X main.version=$(VERSION) -X main.buildTime=$(DATE)" -v -o $(WINDOWS)

$(LINUX):
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -v -ldflags="-s -w -X main.version=$(VERSION) -X main.buildTime=$(DATE)" -o $(LINUX)

$(DARWIN):
	env GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -v -ldflags="-s -w -X main.version=$(VERSION) -X main.buildTime=$(DATE)" -o $(DARWIN)


upgrade:
	go get -u -v
	go mod download
	go mod tidy
	go mod verify

run:
	./portfwd

clean: clean-windows clean-linux clean-darwin

clean-windows:
	rm -f $(WINDOWS)

clean-linux:
	rm -f $(LINUX)

clean-darwin:
	rm -f $(DARWIN)

clean-all: clean
	go clean
	go mod tidy

default: build