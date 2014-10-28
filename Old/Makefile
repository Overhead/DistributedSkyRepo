V=0

build:
	mkdir -p bin
	export GOPATH=${CURDIR}; cd src; make

setup:
	export GOPATH=${CURDIR}; go get github.com/nu7hatch/gouuid; go get code.google.com/p/go.net/websocket; go get github.com/fsouza/go-dockerclient;
	
clean:
	rm -r -f pkg/*
	rm -f bin/*
