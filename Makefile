BINARY := sir-mx-a-lot

all: install

.PHONY: build
build:
	go build -race -o $(BINARY) cmd/sir-mx-a-lot/sir-mx-a-lot.go

.PHONY: install
install:
	go install -i -race -v ./...
