.PHONY: default
ifeq ($(OS),Windows_NT)
EXE=.exe
endif

ifeq ($(shell uname),Linux)
EXE=
endif

# current directory
DIR := ${CURDIR}

# create a udpctl command line utities
default:
	go build -o $(DIR)/bin/udpctl $(DIR)/cmd/main.go && go vet ./...