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
        env CGO_ENABLED=$cgo_enabled GOOS=$os GOARCH=$arch go build -ldflags "$LDFLAGS" -gcflags "$GCFLAGS" -o kcpraw_client_${os}_${arch}${suffix} github.com/ccsexyz/kcpraw/client
        env CGO_ENABLED=$cgo_enabled GOOS=$os GOARCH=$arch go build -ldflags "$LDFLAGS" -gcflags "$GCFLAGS" -o kcpraw_server_${os}_${arch}${suffix} github.com/ccsexyz/kcpraw/server
		if $UPX; then upx -9 kcpraw_client_${os}_${arch}${suffix} kcpraw_server_${os}_${arch}${suffix};fi
		tar -zcf kcpraw-${os}-${arch}-$VERSION.tar.gz kcpraw_client_${os}_${arch}${suffix} kcpraw_server_${os}_${arch}${suffix}
		$MD5 kcpraw-${os}-${arch}-$VERSION.tar.gz
	done
done

# ARM
ARMS=(5 6 7)
for v in ${ARMS[@]}; do
	env CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=$v go build -ldflags "$LDFLAGS" -gcflags "$GCFLAGS" -o kcpraw_client_linux_arm$v  github.com/ccsexyz/kcpraw/client
	env CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=$v go build -ldflags "$LDFLAGS" -gcflags "$GCFLAGS" -o kcpraw_server_linux_arm$v  github.com/ccsexyz/kcpraw/server
done
if $UPX; then upx -9 kcpraw_client_linux_arm* kcpraw_server_linux_arm*;fi
tar -zcf kcpraw-linux-arm-$VERSION.tar.gz kcpraw_client_linux_arm* kcpraw_server_linux_arm*
$MD5 kcpraw-linux-arm-$VERSION.tar.gz

#MIPS32LE
env CGO_ENABLED=0 GOOS=linux GOARCH=mipsle go build -ldflags "$LDFLAGS" -gcflags "$GCFLAGS" -o kcpraw_client_linux_mipsle github.com/ccsexyz/kcpraw/client
env CGO_ENABLED=0 GOOS=linux GOARCH=mipsle go build -ldflags "$LDFLAGS" -gcflags "$GCFLAGS" -o kcpraw_server_linux_mipsle github.com/ccsexyz/kcpraw/server
env CGO_ENABLED=0 GOOS=linux GOARCH=mips go build -ldflags "$LDFLAGS" -gcflags "$GCFLAGS" -o kcpraw_client_linux_mips github.com/ccsexyz/kcpraw/client
env CGO_ENABLED=0 GOOS=linux GOARCH=mips go build -ldflags "$LDFLAGS" -gcflags "$GCFLAGS" -o kcpraw_server_linux_mips github.com/ccsexyz/kcpraw/server

if $UPX; then upx -9 kcpraw_client_linux_mips* kcpraw_server_linux_mips*;fi
tar -zcf kcpraw-linux-mipsle-$VERSION.tar.gz kcpraw_client_linux_mipsle kcpraw_server_linux_mipsle
tar -zcf kcpraw-linux-mips-$VERSION.tar.gz kcpraw_client_linux_mips kcpraw_server_linux_mips
$MD5 kcpraw-linux-mipsle-$VERSION.tar.gz
$MD5 kcpraw-linux-mips-$VERSION.tar.gz
