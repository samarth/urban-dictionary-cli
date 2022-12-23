#!/bin/bash

mkdir -p ./build

platforms=("windows/amd64" "windows/386" "linux/386" "linux/amd64")

for platform in "${platforms[@]}"
do
	platform_split=(${platform//\// })
	GOOS=${platform_split[0]}
	GOARCH=${platform_split[1]}
	output_name="./build/urban-"$GOOS'-'$GOARCH
	if [ $GOOS = "windows" ]; then
		output_name+='.exe'
	fi
	env GOOS=$GOOS GOARCH=$GOARCH go build -o $output_name udcli.go
done

