# poemoon

This tool provides automatic login feature for connect to the network authenticated by dr.com. Tested on drcom client 5.2.0D. Implemented by both drcom and http protocol.

All work done, it will opens a webpage, most browsers will open it on a new tab.

## Dependency

To install, run the following commands in this order:

```sh
go get github.com/josephspurrier/goversioninfo
cd $GOPATH/src/github.com/josephspurrier/goversioninfo/cmd/goversioninfo
go build
mv goversioninfo* $GOROOT/bin
```

## Installation

You should always execute the command below first, so goversioninfo will create a file called resource.syso in the same directory. Then you must run "go build", Go will embed the version information and an optional icon in the executable.

```sh
go generate
```

> **Note:**
> If you use the build.sh to complie the executable, there will be no any Microsoft Windows File Properties/Version Info embed in the executable.

##### for windows
```sh
# for 32 bit
sh build.sh windows 386
# for 64 bit
sh build.sh windows amd64
```

##### for mac os x
```sh
# for 64 bit
sh build.sh darwin amd64
```

##### for linux
```sh
# for 32 bit
sh build.sh linux 386
# for 64 bit
sh build.sh linux amd64
```

That's it!

## Crossplatform support

The package is guaranteed to work on `windows`, `linux` and `darwin`. It also has default support for `freebsd`, `openbsd` and `netbsd` but these three have not been tested yet (that I'm aware of).

## License

It is licensed under the MIT open source license, please see the [LICENSE.md] file for more information.

## Thanks...

panjunwen wrote a nicer version by python. I forked it here [fuzzy_test.py](./fuzzy_test.py), [check it out!](https://github.com/panjunwen/Dr.COM-login).
