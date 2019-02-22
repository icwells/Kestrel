#!/bin/bash

##############################################################################
#	Installs Go scripts for Kestrel
#
#		Requires:	Go 1.11+
##############################################################################

MAIN="kestrel"
FS="github.com/renstrom/fuzzysearch/fuzzy"
GQ="github.com/PuerkitoBio/goquery"
HT="golang.org/x/net/html"
IO="github.com/icwells/go-tools/iotools"
SA="github.com/icwells/go-tools/strarray"
SE="github.com/tebeka/selenium"

# Get install location
SYS=$(ls $GOPATH/pkg | head -1)
PDIR=$GOPATH/pkg/$SYS

echo ""
echo "Preparing Kestrel package..."
echo "GOPATH identified as $GOPATH"
echo ""

# Get dependencies
for I in $FS $GQ $HT $IO $SA $SE ; do
	if [ ! -e "$PDIR/$I.a" ]; then
		echo "Installing $I..."
		go get -u $I
		echo ""
	fi
done

echo "Building main..."
go build -o bin/$MAIN src/*.go

echo ""
echo "Done"
echo ""
