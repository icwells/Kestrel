#!/bin/bash

##############################################################################
# Performs black box tests on kestrel output
#	Requires:	Pytest
##############################################################################
WD=$(pwd)
SRC="$WD/src"

EXTRACTINPUT="test/testInput.csv"
EXPECTED="test/taxonomies.csv"
EXTRACTOUTPUT="test/extracted.csv"
SEARCHOUTPUT="test/searchResults.csv"
REJECTED="test/KestrelRejected.csv"
MISSED="test/KestrelMissed.csv"

whiteBoxTests () {
	echo ""
	echo "Running white box tests..."
	for I in $(ls $SRC); do
		if [ -d $I ]; then
			go test "$SRC/$I*.go"
		fi
	done
}

testExtract () {
	# Extract names and compare output
	kestrel extract -c 0 -i $EXTRACTINPUT -o $EXTRACTOUTPUT
	go test blackBox_test.go --run TestExtract
}

testSearch () {
	# Run search and comapre output
	kestrel search -i $EXTRACTINPUT -o $SEARCHOUTPUT
	go test blackBox_test.go --run TestSearch
}

cleanup () {
	for I in $REJECTED $MISSED $EXTRACTOUTPUT $SEARCHOUTPUT; do
		rm $I
	done
}

blackBoxTests () {
	# Wraps calls to testSearch and testExtract
	./install.sh
	testExtract
	testSearch
	cleanup
}

checkSource () {
	# Runs go fmt/vet on source files
	echo ""
	echo "Running go $1..."
	go $1 "$SRC/main.go"
	#for I in $(ls $SRC); do
	#	if [ -d $I ]; then
	#		go $1 "$SRC/$I*.go"
	#	fi
	#done
}

helpText () {
	echo "Installs Go scripts for Kestrel"
	echo ""
	echo "all			Runs all tests."
	echo "whitebox		Runs white box tests only."
	echo "blackbox		Runs black box tests only."
	echo "help			Prints help text and exits."
	echo "fmt		Runs go fmt on all source files."
	echo "vet		Runs go vet on all source files."
}

if [ $# -eq 0 ]; then
	whiteBoxTests
	cd bin/
	testExtract
	testSearch
	cleanup
elif [ $1 = "whitebox" ]; then
	whiteBoxTests
elif [ $1 = "blackbox" ]; then
	blackBoxTests
elif [ $1 = "all" ]; then
	whiteBoxTests
	blackBoxTests
elif [ $1 = "fmt" ]; then
	checkSource $1
elif [ $1 = "vet" ]; then
	checkSource $1
elif [ $1 = "help" ]; then
	helpText
else
	helpText
fi
echo ""
