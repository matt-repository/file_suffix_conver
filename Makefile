.PHONY: build install clean test integration dep release
GIT_SHA=`git rev-parse --short HEAD || echo`
build:
	@echo "Building fsc..."
	@mkdir -p bin
	@go build -ldflags "-X main.GitSHA=${GIT_SHA}" -o bin/fsc .

install:
	@echo "Installing fsc..."
	@install -c bin/fsc ${GOPATH}/bin/fsc

clean:
	@rm -f bin/*

package:build
	@echo "make app package"
	go-bin-rpm generate --output fsc-0.1.0.rpm --version 0.1.0
