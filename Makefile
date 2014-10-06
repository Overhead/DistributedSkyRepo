V=0

build:
	mkdir -p bin
	export GOPATH=${CURDIR}; cd src; make

setup:
	export GOPATH=${CURDIR}; go get github.com/nu7hatch/gouuid; go get code.google.com/p/go.net/websocket;
	
clean:
	rm -r -f pkg/*
	rm -f bin/*
