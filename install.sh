#!/bin/bash

##############################################################################
#	Installs Go scripts for Kestrel
#
#		Requires:	Go 1.11+
##############################################################################

MAIN="kestrel"
SE="github.com/tebeka/selenium"

installSelenium () {
	# Installs selenium package
	echo "Installing Selenium driver..."
	WD=$(pwd)
	cd "$GOPATH/src/$SE/vendor"
	go get -d ./...
	go run init.go --alsologtostderr
	cd $WD
}

installPackages () {
	echo "Installing dependencies..."
	GOQUERY="github.com/PuerkitoBio/goquery"
	DATAFRAME="github.com/icwells/go-tools/dataframe"
	IOTOOLS="github.com/icwells/go-tools/iotools"
	STRARRAY="github.com/icwells/go-tools/strarray"
	KINGPIN="gopkg.in/alecthomas/kingpin.v2"
	SIMPLESET="github.com/icwells/simpleset"
	FUZZY="github.com/lithammer/fuzzysearch/fuzzy"
	ASPELL="github.com/trustmaster/go-aspell"
	for I in $GOQUERY $DATAFRAME $IOTOOLS $STRARRAY $KINGPIN $SIMPLESET $FUZZY $ASPELL $SE; do
		go get $I
	done
}

installMain () {
	echo "Building main..."
	go build -i -o $GOBIN/$MAIN src/*.go
}

echo ""
echo "Preparing Kestrel package..."
echo "GOPATH identified as $GOPATH"
echo ""

if [ $# -eq 0 ]; then
	installMain
elif [ $1 = "all" ]; then
	installPackages
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
