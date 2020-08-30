#!/usr/bin/bash

DIR=`dirname "$0"`
FILE=bundled.go
BIN=`go env GOPATH`/bin

if [ -z "$1" ]; then
    echo Please specify an image
    exit
fi

if [ ! -f "$1" ]; then
    echo $1 is not a valid file
    exit
fi

cd $DIR

FN=$(basename $1)
BASE=${FN%.*}
$BIN/fyne bundle -package main -append -name "${BASE}Bitmap" $1 > $FILE
