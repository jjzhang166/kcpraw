#!/bin/bash
MD5='md5sum'
unamestr=`uname`
if [[ "$unamestr" == 'Darwin' ]]; then
	MD5='md5'
fi

UPX=false
if hash upx 2>/dev/null; then
	UPX=true
fi

VERSION=`date -u +%Y%m%d`
LDFLAGS="-X main.VERSION=$VERSION -s -w"
GCFLAGS=""

OSES=(linux darwin windows freebsd)
ARCHS=(amd64 386)
for os in ${OSES[@]}; do
	for arch in ${ARCHS[@]}; do
		suffix=""
		if [ "$os" == "windows" ]
		then
			suffix=".exe"
		fi
        cgo_enabled=1
        if [ "$os" == "linux" ]
        then 
            cgo_enabled=0
        fi 
        env CGO_ENABLED=$cgo_enabled GOOS=$os GOARCH=$arch go build -ldflags "$LDFLAGS" -gcflags "$GCFLAGS" -o rawclient_${os}_${arch}${suffix} github.com/ccsexyz/kcptun/client
        env CGO_ENABLED=$cgo_enabled GOOS=$os GOARCH=$arch go build -ldflags "$LDFLAGS" -gcflags "$GCFLAGS" -o rawserver_${os}_${arch}${suffix} github.com/ccsexyz/kcptun/server
		if $UPX; then upx -9 rawclient_${os}_${arch}${suffix} rawserver_${os}_${arch}${suffix};fi
		tar -zcf kcptun-${os}-${arch}-$VERSION.tar.gz rawclient_${os}_${arch}${suffix} rawserver_${os}_${arch}${suffix}
		$MD5 kcptun-${os}-${arch}-$VERSION.tar.gz
	done
done

# ARM
ARMS=(5 6 7)
for v in ${ARMS[@]}; do
	env CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=$v go build -ldflags "$LDFLAGS" -gcflags "$GCFLAGS" -o rawclient_linux_arm$v  github.com/ccsexyz/kcptun/client
	env CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=$v go build -ldflags "$LDFLAGS" -gcflags "$GCFLAGS" -o rawserver_linux_arm$v  github.com/ccsexyz/kcptun/server
done
if $UPX; then upx -9 rawclient_linux_arm* rawserver_linux_arm*;fi
tar -zcf kcptun-linux-arm-$VERSION.tar.gz rawclient_linux_arm* rawserver_linux_arm*
$MD5 kcptun-linux-arm-$VERSION.tar.gz

#MIPS32LE
env CGO_ENABLED=0 GOOS=linux GOARCH=mipsle go build -ldflags "$LDFLAGS" -gcflags "$GCFLAGS" -o rawclient_linux_mipsle github.com/ccsexyz/kcptun/client
env CGO_ENABLED=0 GOOS=linux GOARCH=mipsle go build -ldflags "$LDFLAGS" -gcflags "$GCFLAGS" -o rawserver_linux_mipsle github.com/ccsexyz/kcptun/server
env CGO_ENABLED=0 GOOS=linux GOARCH=mips go build -ldflags "$LDFLAGS" -gcflags "$GCFLAGS" -o rawclient_linux_mips github.com/ccsexyz/kcptun/client
env CGO_ENABLED=0 GOOS=linux GOARCH=mips go build -ldflags "$LDFLAGS" -gcflags "$GCFLAGS" -o rawserver_linux_mips github.com/ccsexyz/kcptun/server

if $UPX; then upx -9 rawclient_linux_mips* rawserver_linux_mips*;fi
tar -zcf kcptun-linux-mipsle-$VERSION.tar.gz rawclient_linux_mipsle rawserver_linux_mipsle
tar -zcf kcptun-linux-mips-$VERSION.tar.gz rawclient_linux_mips rawserver_linux_mips
$MD5 kcptun-linux-mipsle-$VERSION.tar.gz
$MD5 kcptun-linux-mips-$VERSION.tar.gz
