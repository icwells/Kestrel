'''This script contains utility functions for kestrel.'''

import os
import nltk
import gzip
from random import shuffle

cdef str NCBI = "https://eutils.ncbi.nlm.nih.gov/entrez/eutils/"
cdef str EOL= "http://eol.org/api/"
cdef str WIKI = "https://en.wikipedia.org/wiki/"

cdef float L = 0.1

def nameFeatures(name):
	return {"spaces": name.count(" "), "first": name[0], "last": name[-1], "secondlast": name[-2]}

def getNames():
	# Reads common and scientific names in as list of tagged tupples
	cdef list labled = []
	if not os.path.isfile("commonNames.csv.gz"):
		print("\n\t[Error] Common species name file not found. Exiting\.n")
		quit()
	with gzip.open("commonNames.csv.gz", "rb") as f:
		for line in f:
			line = str(line)[2:-3]
			line = line.split(",")
			# Split names into individual words
			labled.append((line[0],"common"))
			labled.append((line[1],"scientific"))
	shuffle(labled)
	return list(labled)

def getSequenceClassifier():
	# Generates classifier for seperating common names and scientific names
	cdef int idx
	cdef list test
	cdef list train
	cdef list features = []
	print("\tGenerating feature classifier...")
	labled = getNames()
	for i in labled:
		# Get features for classifier
		features.append((nameFeatures(i[0]), i[1]))
	idx = int(len(features)*L)
	test, train = features[:idx], features[idx:]
	classifier = nltk.NaiveBayesClassifier.train(train)
	print(("\tFeature classifier evaluation: {:.2%}").format(nltk.classify.accuracy(classifier, test)))
	return classifier

#---------------------------------------------------------------------------------

def apiKeys():
	# Reads in api keys as a dictionary
	cdef str line
	cdef list splt
	keys = {}
	if not os.path.isfile("API.txt"):
		print("\n\t[Error] API key file not found. Exiting.\n")
		quit()
	with open("API.txt", "r") as f:
		for line in f:
			splt = line.split("=")
			if splt[0].strip() == "EOL":
				keys[EOL] = splt[1].strip()
			elif splt[0].strip() == "NCBI":
				keys[NCBI] = splt[1].strip()
	return keys

def speciesList(infile, c, done=[]):
	# Extracts list of query sequences from input file
	cdef int first = 1
	cdef str line
	cdef str delim
	cdef list splt
	q = set()
	print("\tReading input file...")
	with open(infile, "r") as f:
		for line in f:
			if first == 0:
				line = line.strip()
				if delim:
					splt = line.split(delim)
					if len(splt) >= c:
						if splt[c] not in done:
							q.add(splt[c])
				else:
					q.add(line)
			else:
				if "\t" in line:
					delim = "\t"
				elif "," in line:
					delim = ","
				else:
					delim = None
				first = 0
	return list(q)

def checkOutput(outfile, header):
	# Makes output file if needed and reads in any completed queries
	cdef int first = 1
	cdef str line
	done = set()
	if os.path.isfile(outfile):
		print("\tReading previous output...")
		with open(outfile, "r") as output:
			for line in output:
				if first == 0:
					# Save query names
					line = line.strip()
					done.add(line.split(",")[0])
				else:
					# Skip header
					first = 0
	else:
		print("\tGenerating new output file...")
		with open(outfile, "w") as output:
			# Initialize file and write header
			output.write(header)
	return list(done)
