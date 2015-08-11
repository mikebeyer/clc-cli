VERSION=0.7.0

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
	gox -output="clc" -osarch="linux/amd64"
	tar czf clc_linux.$(VERSION).tar.gz clc
	rm clc
	gox -output="clc" -osarch="darwin/amd64"
	tar czf clc_darwin.$(VERSION).tar.gz clc
	rm clc
	gox -output="clc" -osarch="windows/amd64"
	tar czf clc_windows.$(VERSION).tar.gz clc.exe
	rm clc.exe
