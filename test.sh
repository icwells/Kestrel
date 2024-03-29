#!/bin/bash

##############################################################################
# Performs black box tests on kestrel output
#	Requires:	Pytest
##############################################################################
USER=""
PW=""
WD=$(pwd)
SRC="$WD/src"
TEST="$WD/test"

EXTRACTINPUT="$TEST/testInput.csv"
EXPECTED="$WD/utils/corpus.csv.gz"
SEARCHOUTPUT="$TEST/searchResults.csv"
REJECTED="$TEST/KestrelRejected.csv"
MISSED="$TEST/KestrelMissed.csv"

SEARCHTAXA="$SRC/searchtaxa/*.go"
TAXONOMY="$SRC/taxonomy/*.go"
TERMS="$SRC/terms/*.go"

getUser () {
	# Reads mysql user name and password from command line
	read -p "Enter MySQL username: " USER
	echo -n "Enter MySQL password: "
	read -s PW
	echo ""
}

whiteBoxTests () {
	echo ""
	echo "Running white box tests..."
	go test $SEARCHTAXA
	go test $TAXONOMY
	go test $TERMS
}

testSearch () {
	# Run search and compare output
	getUser
	cleanup
	go run src/main.go search -u $USER --password $PW -c 1 -i $EXTRACTINPUT -o $SEARCHOUTPUT
	cd $TEST
	go test blackBox_test.go --run TestSearch
	cleanup
}

fullSearch () {
	# Runs large scale black box tests
	getUser
	cd $TEST
	cleanup
	go run accuracyTest.go --user $USER --password $PW
}

cleanup () {
	for I in $REJECTED $MISSED $EXTRACTOUTPUT $SEARCHOUTPUT; do
		if [ -f $I ]; then
			rm $I
		fi
	done
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
	go $1 "$TEST/accuracyTest.go"
}

helpText () {
	echo "Installs Go scripts for Kestrel"
	echo ""
	echo "all			Runs all tests."
	echo "whitebox		Runs white box tests only."
	echo "blackbox		Runs black box tests only."
	echo "full			Performs large scale black box test."
	echo "fmt		Runs go fmt on all source files."
	echo "vet		Runs go vet on all source files."
	echo "help			Prints help text and exits."
}

if [ $# -eq 0 ]; then
	helpText
elif [ $1 = "whitebox" ]; then
	whiteBoxTests
elif [ $1 = "blackbox" ]; then
	testSearch
elif [ $1 = "all" ]; then
	whiteBoxTests
	testSearch
elif [ $1 = "full" ]; then
	fullSearch
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
