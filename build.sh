#!/bin/sh

#####################################################################
# usage:
# sh build.sh 构建默认的windows 32位程序
# sh build.sh darwin(或linux), 构建指定平台的64位程序

# examples:
# sh build.sh darwin amd64 构建MacOS版本的64位程序
# sh build.sh linux amd64 构建linux版本的64位程序
#####################################################################

source /etc/profile

OS="$1"
ARCH="$2"

if [ -n "$OS" ];then
   echo "use defined GOOS: "${OS}
else
   echo "use default GOOS: windows"
   OS=windows
   echo "use default GOOS: 386"
   ARCH=386
fi

echo "start building with GOOS: "${OS}", GOARCH: "${ARCH}

if [ ${OS} == "windows" ];then
   SUFFIX=".exe"
else
   SUFFIX=""
fi

export GOOS=${OS}
export GOARCH=${ARCH}


release_dir="poemoon"
revision=`git describe --long --dirty`


mkdir -p ./${release_dir}
rm -rf ./${release_dir}/*


flags="-X main.buildstamp=`date -u '+%Y-%m-%d_%I:%M:%S%p'` -X main.githash=`git describe --long --dirty --abbrev=14` -X 'main.goversion=`go version`'"
echo ${flags}
go build -ldflags "$flags" -x -o ${release_dir}/drcom_hp${SUFFIX} cmd/drcom_hp/main.go
go build -ldflags "$flags" -x -o ${release_dir}/poemoon${SUFFIX} cmd/poemoon/main.go

cp -r config/* ./${release_dir}/
cp ./README.md ./${release_dir}/


echo "finish building with GOOS: "${OS}", GOARCH: "${ARCH}

rm -rf poemoon.tar.gz
tar zcvf poemoon.tar.gz poemoon
