#!/bin/bash

##############################################################################
#	Installs Go scripts for Kestrel
#
#		Requires:	Go 1.11+
##############################################################################

DIR="databases"
MAIN="kestrel"
PW=""
SE="github.com/tebeka/selenium"
USER=""

getUser () {
	# Reads mysql user name and password from command line
	read -p "Enter MySQL username: " USER
	echo -n "Enter MySQL password: "
	read -s PW
	echo ""
}

downloadDatabases () {
	ITIS="https://www.itis.gov/downloads/itisMySQLBulk.zip"
	getUser
	mkdir $DIR
	cd $DIR
	echo "Downloading databases..."
	wget $ITIS
	echo "Extracting files..."
	unzip itisMySQL*
	echo "Uploading ITIS tables to MySQL..."
	mv itisMySQL*/* .
	mysql -u$USER -p$PW < CreateDB.sql
	cd ../
	rm -r $DIR
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

installTensorFlow () {
	TF="libtensorflow-cpu-linux-x86_64-2.6.0.tar.gz"
	LINK="https://storage.googleapis.com/tensorflow/libtensorflow/$TF"
	wget $LINK
	sudo tar -C /usr/local -xzf $TF
	sudo ldconfig
	rm $TF
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
	TFGO="github.com/galeone/tfgo"
	for I in  $ASPELL $DATAFRAME $DBIO $FUZZY $GOQUERY $IOTOOLS $KINGPIN $SE $SIMPLESET $STRARRAY $TFGO; do
		go get $I
	done
	# Install TensorFlow Go bindings
	go env -w GONOSUMDB="github.com/galeone/tensorflow"
}

installMain () {
	echo "Building main..."
	go build -i -o $GOBIN/$MAIN src/*.go
}

helpText () {
	echo "Installs Go scripts for Kestrel"
	echo ""
	echo "all	Installs all depenencies, including TensorFlow, Selenium package, and drivers."
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
	installTensorFlow
	installSelenium
	installPackages
	installMain
elif [ $1 = "download" ]; then
	downloadDatabases
elif [ $1 = "help" ]; then
	helpText
fi

echo ""
echo "Done"
echo ""
