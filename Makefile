VERSION=0.1

.PHONY : build deps clean
build:
	godep go build -o clc
deps:
	go get github.com/tools/godep
	godep restore
clean:
	rm clc
