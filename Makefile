VERSION ?= 0.0.1  
NAME ?= "Anonymous Chatting Platform"  
AUTHOR ?= "Timileyin Farayola"   
  
.PHONY: build run start fp p

build:  
	@cd ./cmd && go build -o ../chat

run:
	@go run ./cmd/main.go

start:
	make && ./chat 2>&1 | tee chat_app_logs.txt

fp:
	git push -f

p:
	git push


DEFAULT: build
