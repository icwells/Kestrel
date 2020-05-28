#!/bin/bash

##############################################################################
# Performs black box tests on kestrel output
#	Requires:	Pytest
##############################################################################
WD=$(pwd)
SRC="$WD/src"
TEST="$WD/test"

EXTRACTINPUT="$TEST/testInput.csv"
EXPECTED="$TEST/taxonomies.csv"
SEARCHOUTPUT="$TEST/searchResults.csv"
REJECTED="$TEST/KestrelRejected.csv"
MISSED="$TEST/KestrelMissed.csv"

SEARCHTAXA="$SRC/searchtaxa/*.go"
TAXONOMY="$SRC/taxonomy/*.go"
TERMS="$SRC/terms/*.go"

whiteBoxTests () {
	echo ""
	echo "Running white box tests..."
	go test $SEARCHTAXA
	go test $TAXONOMY
	go test $TERMS
}

testSearch () {
	# Run search and comapre output
	cd $TEST
	kestrel search -i $EXTRACTINPUT -o $SEARCHOUTPUT
	go test blackBox_test.go --run TestSearch
}

cleanup () {
	for I in $REJECTED $MISSED $EXTRACTOUTPUT $SEARCHOUTPUT; do
		if [ -f $I ]; then
			rm $I
		fi
	done
}

blackBoxTests () {
	# Wraps calls to testSearch and testExtract
	./install.sh
	testSearch
	cleanup
}

checkSource () {
	# Runs go fmt/vet on source files
	echo ""
	echo "Running go $1..."
	go $1 "$SRC/main.go"
	go $1 $SEARCHTAXA
	go $1 $TAXONOMY
	go $1 $TERMS
	go $1 "$TEST/blackBox_test.go"
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
	helpText
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
