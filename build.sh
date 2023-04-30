#!/bin/bash
BINDIR=$(dirname "$0")/bin
BUILDOS=$1
CDIR=${PWD##*/} # to assign to a variable

echo Building "$CDIR"

if [ -z "$1" ]; then
    BUILDOS="linux"
    echo "Statically Building in $BINDIR for $BUILDOS"
    GOOS=$BUILDOS CGO_ENABLED=0 go build -o "$BINDIR"/"$CDIR"

    BUILDOS="windows"
    echo "Statically Building in $BINDIR for $BUILDOS"
    GOOS=$BUILDOS CGO_ENABLED=0 go build -o "$BINDIR"/"$CDIR".exe

    BUILDOS="darwin"
    echo "Statically Building in $BINDIR for $BUILDOS"
    GOOS=$BUILDOS CGO_ENABLED=0 go build -o "$BINDIR"/"$CDIR"_darwin

    BUILDOS="darwin"
    echo "Statically Building in $BINDIR for $BUILDOS"
    GOOS=$BUILDOS GOARCH=arm64 CGO_ENABLED=0 go build -o "$BINDIR"/"$CDIR"_darwin_arm64
else
    echo "Statically Building in $BINDIR for $BUILDOS"
    GOOS=$BUILDOS CGO_ENABLED=0 go build -o "$BINDIR"/"$CDIR"
fi
