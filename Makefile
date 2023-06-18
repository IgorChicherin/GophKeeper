PROJECTNAME=gophkeeper

PATH_TO_DARWIN_AMD64=bin/macos-amd64
PATH_TO_LINUX_AMD64=bin/linux-amd64
PATH_TO_WINDOWS_AMD64=bin/windows-amd64

PATH_FROM_DARWIN_AMD64=fyne-cross/dist/darwin-amd64
PATH_FROM_LINUX_AMD64=fyne-cross/dist/linux-amd64
PATH_FROM_WINDOWS_AMD64=fyne-cross/dist/windows-amd64

TEAM_ID=ashzak

default: help

create-download-dir:
	mkdir -p download
	mkdir -p $(PATH_TO_DARWIN_AMD64)
	mkdir -p $(PATH_TO_LINUX_AMD64)
	mkdir -p $(PATH_TO_WINDOWS_AMD64)


build-darwin-amd64:
	fyne-cross darwin -arch amd64 -macosx-sdk-path="/app/sdk/MacOSX.sdk" -no-cache -debug -name $(PROJECTNAME) -app-id $(TEAM_ID).$(PROJECTNAME) ./cmd/client

build-linux-amd64:
	fyne-cross linux -arch amd64 -name $(PROJECTNAME) ./cmd/client

build-windows-amd64:
	fyne-cross windows -arch amd64 -name $(PROJECTNAME).exe -app-id $(TEAM_ID).$(PROJECTNAME) ./cmd/client

pack-darwin-amd64:
	mv $(PATH_FROM_DARWIN_AMD64)/$(PROJECTNAME).app $(PATH_TO_DARWIN_AMD64)/$(PROJECTNAME).app

pack-linux-amd64:
	mv $(PATH_FROM_LINUX_AMD64)/$(PROJECTNAME).tar.xz $(PATH_TO_LINUX_AMD64)/$(PROJECTNAME).tar.xz

pack-windows-amd64:
	mv $(PATH_FROM_WINDOWS_AMD64)/$(PROJECTNAME).exe.zip $(PATH_TO_WINDOWS_AMD64)/$(PROJECTNAME).zip

remove-binares:
	rm -rf fyne-cross

remove-archives:
	rm -rf bin

clean: remove-archives remove-binares

build: build-linux-amd64 build-windows-amd64 build-darwin-amd64

pack:  pack-linux-amd64 pack-windows-amd64 pack-darwin-amd64

prepare-release: create-download-dir build pack remove-binares

run-server:
	docker-compose -f ./deployments/docker-compose.yml up --build db adminer gophkeeper-server

build-clients:
	docker-compose -f ./deployments/docker-compose.yml up --build client-builder

build-all:
	docker-compose -f ./deployments/docker-compose.yml up --build

install:
	fyne install

help:
	@echo "make <command>"
	@echo "The commands are:"
	@echo "install                  install on your system"
	@echo "run-server               run server stack in docker"
	@echo "prepare-release          build for each os, pack, mv and clean"
	@echo "create-download-dir      create download and subdirs"
	@echo "pack                     pack and mv for each os"
	@echo "build                    build for each os"
	@echo "build-darwin-amd64       build for macos intel"
	@echo "build-linux-amd64        build for linux x86-64"
	@echo "build-windows-amd64      build for windows x86-64"
	@echo "pack-darwin-arm64        pack and mv for macos arm"
	@echo "pack-darwin-amd64        pack and mv for macos intel"
	@echo "build-linux-amd64        pack and mv for linux x86-64"
	@echo "build-windows-amd64      pack and mv for windows x86-64"
	@echo "remove-binares           remove dir fyne-cross"
	@echo "remove-archives          remove dir download"
	@echo "clean                    clean all"

all: help