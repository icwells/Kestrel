#!/bin/bash

##############################################################################
#	Installs Go scripts for Kestrel
#
#		Requires:	Go 1.11+
##############################################################################

MAIN="kestrel"
SE="github.com/tebeka/selenium"

# Get install location
SYS=$(ls $GOPATH/pkg | head -1)
PDIR=$GOPATH/pkg/$SYS

installSelenium () {
	# Installs selenium package
	WD=$(pwd)
	installPackage $SE
	cd $GOPATH/$SE/vendor
	go get -d ./...
	go run init.go --alsologtostderr
	cd $WD
}

installMain () {
	echo "Building main..."
	go build -i -o bin/$MAIN src/*.go
}

echo ""
echo "Preparing Kestrel package..."
echo "GOPATH identified as $GOPATH"
echo ""

if [ $# -eq 0 ]; then
	installMain
elif [ $1 = "all" ]; then
	installSelenium
	installMain
elif [ $1 = "help" ]; then
	echo "Installs Go scripts for Kestrel"
	echo ""
	echo "all	Installs all depenencies, including selenium package and drivers."
	echo "help	Prints help text and exits."
	echo ""
fi

echo ""
echo "Done"
echo ""
