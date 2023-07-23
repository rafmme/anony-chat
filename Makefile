VERSION ?= 0.0.1  
NAME ?= "Anonymous Chatting Platform"  
AUTHOR ?= "Timileyin Farayola"   
  
.PHONY: build run fp p

build:  
	@cd ./cmd && go build -o ../chat

run:
	@go run ./cmd/main.go

fp:
	git push -f

p:
	git push


DEFAULT: build
