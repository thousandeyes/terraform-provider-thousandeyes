GOPATH?=$(shell go env GOPATH)
GO111MODULE=auto

build:
	go build -o terraform-provider-thousandeyes
