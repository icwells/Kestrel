#!/bin/bash

##############################################################################
# Performs black box tests on kestrel output
#	Requires:	Pytest
##############################################################################

EXTRACTINPUT="test/testInput.csv"
EXPECTED="test/taxonomies.csv"
EXTRACTOUTPUT="test/extracted.csv"
SEARCHOUTPUT="test/searchResults.csv"
REJECTED="test/KestrelRejected.csv"
MISSED="test/KestrelNoMatch.csv"

cd bin/

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

testExtract
testSearch
cleanup
