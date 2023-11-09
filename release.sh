#!/bin/bash

# look at this first
# https://github.com/marketplace/actions/go-release-binaries
set -e

tag="${1}"

if [ "${tag}" == "" ]; then
  echo "valid tag value required ${tag}}"
  exit 1
fi

rm -rf dist
GOOS=darwin GOARCH=amd64 go build -o "dist/darwin-x86_64-${tag}"
GOOS=darwin GOARCH=arm64 go build -o "dist/darwin-arm64-${tag}"
GOOS=linux GOARCH=386 go build -o "dist/linux-i386-${tag}"
GOOS=linux GOARCH=amd64 go build -o "dist/linux-x86_64-${tag}"
GOOS=windows GOARCH=386 go build -o "dist/windows-i386-${tag}"
GOOS=windows GOARCH=amd64 go build -o "dist/windows-x86_64-${tag}"

gh release create "$tag" ./dist/* --title="${tag}" --notes "${tag}"
