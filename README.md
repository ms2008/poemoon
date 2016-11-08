# poemoon

This tool provides automatic login feature for connect to the network authenticated by dr.com.

All work done, it will opens a webpage, most browsers will open it on a new tab.

## Installation

### for windows
```sh
# for 32 bit
sh build.sh windows 386
# for 64 bit
sh build.sh windows amd64
```

### for linux
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
