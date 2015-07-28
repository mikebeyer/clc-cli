VERSION=0.1

.PHONY : build deps clean
build:
	godep go build -o clc
deps:
	go get github.com/tools/godep
	godep restore
clean:
	-rm clc
	-rm *.tar.gz
	-rm *.bin
	-rm *.exe
release:
	gox -output="clc_{{.OS}}.bin" -osarch="linux/amd64 darwin/amd64 windows/amd64"
	tar czf clc_linux.$(VERSION).tar.gz clc_linux.bin
	tar czf clc_darwin.$(VERSION).tar.gz clc_darwin.bin
	tar czf clc_windows.$(VERSION).tar.gz clc_windows.bin.exe
	-rm *.bin
	-rm *.exe
