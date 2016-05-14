#!/bin/sh

VERSION=`date -u +%Y%m%d`
LDFLAGS="-X main.VERSION=$VERSION -s -w"

OSES=(linux darwin windows)
ARCHS=(amd64 386)
for os in ${OSES[@]}; do
	for arch in ${ARCHS[@]}; do
		suffix=""
		if [ "$os" == "windows" ]
		then
			suffix=".exe"
		fi
		env GOOS=$os GOARCH=$arch go build -ldflags "$LDFLAGS" -o client_${os}_${arch}${suffix} github.com/xtaci/kcptun/client
		env GOOS=$os GOARCH=$arch go build -ldflags "$LDFLAGS" -o server_${os}_${arch}${suffix} github.com/xtaci/kcptun/server
		tar -zcf kcptun-${os}-${arch}-$VERSION.tar.gz client_${os}_${arch}${suffix} server_${os}_${arch}${suffix}
		md5 kcptun-${os}-${arch}-$VERSION.tar.gz
	done
done

# ARM
ARMS=(5 6 7)
for v in ${ARMS[@]}; do
	env GOOS=linux GOARCH=arm GOARM=$v go build -ldflags "$LDFLAGS" -o client_linux_arm$v  github.com/xtaci/kcptun/client
	env GOOS=linux GOARCH=arm GOARM=$v go build -ldflags "$LDFLAGS" -o server_linux_arm$v  github.com/xtaci/kcptun/server
done
tar -zcf kcptun-linux-arm-$VERSION.tar.gz client_linux_arm* server_linux_arm*
md5 kcptun-linux-arm-$VERSION.tar.gz
