#!/bin/bash

##############################################################################
#	Installs Go scripts for Kestrel
#
#		Requires:	Go 1.11+
##############################################################################

MAIN="src/*.go"
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
for I in $DR $IO $SA $SE ; do
	if [ ! -e "$PDIR/$I.a" ]; then
		echo "Installing $I..."
		go get -u $I
		echo ""
	fi
done

echo "Building main..."
go build -o bin/$MAIN src/$MAIN/*.go

echo ""
echo "Done"
echo ""
