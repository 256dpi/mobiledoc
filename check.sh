#!/usr/bin/env bash

go fmt ./...
go vet ./...
golint ./...
go test ./...
