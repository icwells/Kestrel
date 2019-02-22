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

installPackage () {
	# Installs go package if it is not present in src directory
	if [ ! -e "$PDIR/$1.a" ]; then
		echo "Installing $1..."
		go get -u $1
		echo ""
	fi
}

installSelenium () {
	# Installs selenium package
	WD=$(pwd)
	installPackage $SE
	cd $GOPATH/$SE/vendor
	go get -d ./...
	go run init.go --alsologtostderr
	cd $WD
}

installDependencies () {
	# Get dependencies
	for I in $FS $GQ $HT $IO $SA $ST ; do
		installPackage $I
	done
}

installMain () {
	echo "Building main..."
	go build -o bin/$MAIN src/*.go
}

echo ""
echo "Preparing Kestrel package..."
echo "GOPATH identified as $GOPATH"
echo ""

if [ $# -eq 0 ]; then
	installMain
elif [ $1 = "all" ]; then
	installSelenium
	installDependencies
	installMain
elif [ $1 = "all" ]; then
	echo "Installs Go scripts for Kestrel"
	echo ""
	echo "all	Installs all depenencies, including selenium package and drivers."
	echo "help	Prints help text and exits."
	echo ""
fi

echo ""
echo "Done"
echo ""
