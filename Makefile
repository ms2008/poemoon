.PHONY: build dep clean

export GO15VENDOREXPERIMENT=1

all: clean build

build:
	@sh build.sh linux amd64

mac: clean
	@sh build.sh darwin amd64

win: clean
	@sh build.sh windows amd64

dep:
	@dep ensure

clean:
	@rm -rf poemoon*
