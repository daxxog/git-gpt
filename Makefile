SHELL := /bin/bash


.PHONY: build
build: git-gpt


.PHONY: help
help:
	@printf "available targets -->\n\n"
	@cat Makefile | grep ".PHONY" | grep -v ".PHONY: _" | sed 's/.PHONY: //g'


git-gpt: main.go message.go config.go
	go build .


.PHONY: install
install: git-gpt
	cp git-gpt /usr/local/bin
