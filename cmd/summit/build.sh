#!/bin/bash
VERSION=$(git describe --tags --abbrev=0)
echo Building Summit version $VERSION
go build -ldflags "-X github.com/RhykerWells/Summit/common.VERSION=$VERSION"
