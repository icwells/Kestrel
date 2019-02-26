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
MISSED="test/KestrelNoMatch.csv"

cd bin/

whiteBoxTests () {
	echo ""
	echo "Running white box tests..."
	go test $SRC
}

testExtract () {
	# Extract names and compare output
	python kestrel.py --extract -c 0 -i $EXTRACTINPUT -o $EXTRACTOUTPUT
	pytest test_kestrel.py::test_extract
}

testSearch () {
	# Run search and comapre output
	python kestrel.py -t 4 -i $EXTRACTINPUT -o $SEARCHOUTPUT
	pytest test_kestrel.py::test_search
}

cleanup () {
	for I in $REJECTED $MISSED $EXTRACTOUTPUT $SEARCHOUTPUT; do
		rm $I
	done
}

#testExtract
#testSearch
#cleanup

if [ $# -eq 0 ]; then
	whiteBoxTests
elif [ $1 = "whitebox" ]; then
	whiteBoxTests
elif [ $1 = "all" ]; then
	whiteBoxTests
elif [ $1 = "help" ]; then
	echo "Installs Go scripts for Kestrel"
	echo ""
	echo "all				Runs all tests."
	echo "whiteBoxTests		Runs white box tests only."
	echo "help				Prints help text and exits."
	echo ""
fi
