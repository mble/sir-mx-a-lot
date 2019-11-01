BINARY := sir-mx-a-lot
PWD := $(pwd)
all:
	go build -race -o $(BINARY) cmd/sir-mx-a-lot/sir-mx-a-lot.go
