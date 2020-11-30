#!/bin/bash
go get
echo "Running short tests"
go test ./... -count=1
echo "Running long tests"
go test ./... -tags=long -count=1