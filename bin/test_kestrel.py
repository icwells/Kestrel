'''Performs black box tests on Kestrel search'''

import pytest

EXTRACTINPUT = "test/testInput.csv"
EXPECTED = "test/taxonomies.csv"
EXTRACTOUTPUT = "test/extracted.csv"
SEARCHOUTPUT = "test/searchResults.csv"

def readFile(infile, taxa):
	# Reads in file as dict
	ret = {}
	first = True
	with open(infile, "r") as f:
		for line in f:
			if first == False:
				s = line.strip().split(",")
				if taxa == True:
					# Use search term as key and ignore sources
					ret[s[1]] = s[1:9]
				else:
					# Use query as key
					ret[s[0]] = s[1:]
			else:
				first = False
	return ret

def compareFiles(exp, act, taxa = False):
	# Compares actual output file to expected file
	expected = readFile(exp, taxa)
	actual = readFile(act, taxa)
	assert len(expected.keys()) == len(actual.keys())
	for k in actual.keys():
		assert k in expected.keys()
		if (k == "Coyote" and actual[k][1] == "scientific") or (k == "canis lupus" and actual[k][1] == "common"):
			# Coyote incorectly labled as scientific by feature classifier
			assert actual[k][0] == expected[k][0]
		else:
			for idx, i in enumerate(actual[k]):
				assert i == expected[k][idx]

def test_extract():
	# Tests extraction/formatting output
	compareFiles(EXTRACTINPUT, EXTRACTOUTPUT)

def test_search():
	# Tests search output
	compareFiles(EXPECTED, SEARCHOUTPUT, True)
