#!/bin/bash

##############################################################################
# Performs black box tests on kestrel output
#	Requires:	Pytest
##############################################################################
WD=$(pwd)
SRC="$WD/src/*.go"

EXTRACTINPUT="test/testInput.csv"
EXPECTED="test/taxonomies.csv"
EXTRACTOUTPUT="test/extracted.csv"
SEARCHOUTPUT="test/searchResults.csv"
REJECTED="test/KestrelRejected.csv"
MISSED="test/KestrelMissed.csv"

whiteBoxTests () {
	echo ""
	echo "Running white box tests..."
	go test $SRC
}

testExtract () {
	# Extract names and compare output
	./kestrel extract -c 0 -i $EXTRACTINPUT -o $EXTRACTOUTPUT
	go test blackBox_test.go --run TestExtract
}

testSearch () {
	# Run search and comapre output
	./kestrel search -i $EXTRACTINPUT -o $SEARCHOUTPUT
	go test blackBox_test.go --run TestSearch
}

cleanup () {
	for I in $REJECTED $MISSED $EXTRACTOUTPUT $SEARCHOUTPUT; do
		rm $I
	done
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
	cd bin/
	testExtract
	testSearch
	cleanup
elif [ $1 = "all" ]; then
	whiteBoxTests
	cd bin/
	testExtract
	testSearch
	cleanup
elif [ $1 = "help" ]; then
	echo "Installs Go scripts for Kestrel"
	echo ""
	echo "all			Runs all tests."
	echo "whitebox		Runs white box tests only."
	echo "blackbox		Runs black box tests only."
	echo "help			Prints help text and exits."
	echo ""
fi
