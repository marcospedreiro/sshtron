# Makefile for sshtron

VERSION ?= 1.0.0
SERVICE_NAME := sshtron

all: clean build run

build:
	CGO_ENABLED=0 go build -ldflags '-X github.com/marcospedreiro/sshtron/version.VERSION=$(VERSION)' -o bin/$(SERVICE_NAME)

run:
	./bin/$(SERVICE_NAME)

clean:
	rm -f bin/*

.PHONY: build run clean
