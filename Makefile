V=0

build:
	mkdir -p bin
	export GOPATH=${CURDIR}; cd src; make

setup:
	export GOPATH=${CURDIR}; go get github.com/nu7hatch/gouuid;
	
clean:
	rm -r -f pkg/*
	rm -f bin/*
