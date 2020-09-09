#!/bin/bash

##############################################################################
#	Installs Go scripts for Kestrel
#
#		Requires:	Go 1.11+
##############################################################################

MAIN="kestrel"
SE="github.com/tebeka/selenium"

downloadDatabases () {
	GBIF="https://hosted-datasets.gbif.org/datasets/backbone/backbone-current-simple.txt.gz"
	ITIS="https://www.itis.gov/downloads/itisMySQLBulk.zip"
	NCBI="https://ftp.ncbi.nlm.nih.gov/pub/taxonomy/taxdump.tar.gz"
	mkdir databases
	cd databases/
	for I in $GBIF $ITIS $NCBI; do
		wget $I
	done
}

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
	ASPELL="github.com/trustmaster/go-aspell"
	DATAFRAME="github.com/icwells/go-tools/dataframe"
	DBIO="github.com/icwells/dbIO"
	FUZZY="github.com/lithammer/fuzzysearch/fuzzy"
	GOQUERY="github.com/PuerkitoBio/goquery"
	IOTOOLS="github.com/icwells/go-tools/iotools"
	KINGPIN="gopkg.in/alecthomas/kingpin.v2"
	SIMPLESET="github.com/icwells/simpleset"
	STRARRAY="github.com/icwells/go-tools/strarray"
	for I in  $ASPELL $DATAFRAME $DBIO $FUZZY $GOQUERY $IOTOOLS $KINGPIN $SE $SIMPLESET $STRARRAY; do
		go get $I
	done
}

installMain () {
	echo "Building main..."
	go build -i -o $GOBIN/$MAIN src/*.go
}

helpText () {
	echo "Installs Go scripts for Kestrel"
	echo ""
	echo "all	Installs all depenencies, including selenium package and drivers."
	echo "download Downloads taxonomy databases (takes several hours)"
	echo "help	Prints help text and exits."
	echo ""
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
elif [ $i = "download" ]; then
	downloadDatabases
elif [ $1 = "help" ]; then
	helpText
fi

echo ""
echo "Done"
echo ""
