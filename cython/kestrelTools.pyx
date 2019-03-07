'''This script contains utility functions for kestrel.'''

import os
import nltk
import gzip
from re import sub
from sys import stdout
from string import punctuation, digits
from random import shuffle

cdef str COMMONNAMES = "test/commonNames.csv.gz"
cdef str NCBI = "https://eutils.ncbi.nlm.nih.gov/entrez/eutils/"
cdef str EOL= "http://eol.org/api/"
cdef str WIKI = "https://en.wikipedia.org/wiki/"
cdef str IUCN = "http://apiv3.iucnredlist.org/api/v3/"

cdef float L = 0.1

def nameFeatures(name):
	return {"spaces": name.count(" "), "first": name[0], "last": name[-1], "secondlast": name[-2]}

def getNames():
	# Reads common and scientific names in as list of tagged tupples
	cdef list labled = []
	if not os.path.isfile(COMMONNAMES):
		print("\n\t[Error] Common species name file not found. Exiting\.n")
		quit()
	with gzip.open(COMMONNAMES, "rb") as f:
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
	print("\n\tGenerating feature classifier...")
	labled = getNames()
	for i in labled:
		# Get features for classifier
		features.append((nameFeatures(i[0]), i[1]))
	idx = int(len(features)*L)
	test, train = features[:idx], features[idx:]
	classifier = nltk.NaiveBayesClassifier.train(train)
	print(("\tFeature classifier evaluation: {:.2%}").format(nltk.classify.accuracy(classifier, test)))
	return classifier

#----------------------------merging-------------------------------------------

def getDelim(line):
	# Returns delimiter
	cdef str i
	for i in ["\t", ",", " "]:
		if i in line:
			return i
	print("\n\t[Error] Cannot determine delimeter. Check file formatting. Exiting.\n")
	quit()

def mergeTaxonomy(infile, outfile, col, taxa):
	# Appends data from taxa to matches in infile
	cdef int first = 1
	cdef str delim
	cdef str line
	cdef list spli
	cdef list na
	cdef list row = []
	cdef int count = 0
	cdef int total = 0
	na = ["NA","NA","NA","NA","NA","NA","NA"]
	print("\tMerging files...")
	with open(outfile, "w") as out:
		with open(infile, "r") as f:
			for line in f:
				total += 1
				if first == 0:
					spli = line.strip().split(delim)
					if len(spli) >= col:
						n = spli[col]
						if n in taxa.keys():
							row = spli + taxa[n]
						else:
							row = spli + na
						if row:
							out.write(",".join(row) + "\n")
							count += 1
				else:
					delim = getDelim(line)
					out.write(line.strip() + ",Kingdom,Phylum,Class,Order,Family,Genus,ScientificName\n")
					first = 0
	print(("\tFound taxonomies for {} of {} entries.\n").format(count, total))

def getTaxa(infile):
	# Reads in taxonomy dictionary
	cdef int first = 1
	cdef list spli
	cdef str line
	taxa = {}
	species = set()
	print("\tReading taxonomies...")
	with open(infile, "r") as f:
		for line in f:
			if first == 0:
				spli = line.split(",")
				# Query name: [taxonomy] (drops search term and urls)
				if len(spli) >= 9:
					taxa[spli[0]] = spli[2:9]
					species.add(spli[8])
			else:
				first = 0
	print(("\tFound {} taxonomies with {} unique species.").format(len(taxa.keys()), len(species)))
	return taxa

#----------------------------i/o----------------------------------------------

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
			elif splt[0].strip() == "IUCN":
				if len(splt) >= 2:
					# Skip if IUCN key is not present
					keys[IUCN] = splt[1].strip()
	if EOL not in keys.keys() or NCBI not in keys.keys():
		print("\n\t[Error] API keys required for EOL and NCBI. Exiting.")
		quit()
	return keys

def termList(infile, done=[]):
	# Reads formatted species names in as a dict
	cdef str line
	cdef list splt
	cdef int first = 1
	terms = {}
	t = []
	with open(infile, "r") as f:
		for line in f:
			if first == 0:
				splt = line.strip().split(",")
				if len(splt) >= 3:
					if splt[0] not in done:
						if splt[1] in terms.keys():
							terms[splt[1]].append(splt[0])
						else:
							# {search term: [term, type, query name]}
							terms[splt[1]] = [splt[1], splt[2], splt[0]]
			else:
				first = 0
	# Convert to list
	for i in terms.keys():
		if len(terms[i]) >= 3:
			t.append(terms[i])
	return t

def speciesList(infile, c, done=[]):
	# Extracts list of query sequences from input file
	cdef int first = 1
	cdef str line
	cdef str delim
	cdef list splt
	cdef int length
	cdef int total = 0
	q = set()
	print("\n\tReading input file...")
	with open(infile, "r", errors="ignore") as f:
		for line in f:
			if first == 0:
				total += 1
				line = line.strip()
				if delim:
					splt = line.split(delim)
					if len(splt) >= length:
						if len(splt[c]) >=1 and splt[c] != "NA":
							# Skip improperly formatted lines
							q.add(splt[c])
				else:
					q.add(line)
			else:
				if "\t" in line:
					delim = "\t"
					length = len(line.split(delim))
				elif "," in line:
					delim = ","
					length = len(line.split(delim))
				else:
					delim = None
					length = 1
				first = 0
	print(("\tFound {} unique entries from {} total entries.").format(len(q), total))
	return list(q)

def checkOutput(outfile, header=""):
	# Makes output file if needed and reads in any completed queries
	cdef int first = 1
	cdef str line
	done = []
	if os.path.isfile(outfile):
		print(("\tReading previous output from {}").format(outfile))
		with open(outfile, "r") as output:
			for line in output:
				if first == 0:
					# Save query names
					line = line.strip()
					done.append(line.split(",")[0])
				else:
					# Skip header
					first = 0
	else:
		print("\tGenerating new output file...")
		with open(outfile, "w") as output:
			# Initialize file and write header
			if header:
				output.write(header)
	return done

def writeResults(outfile, line):
	# Writes match to outfile and no match to misses
	with open(outfile, "a") as output:
		output.write(line)

#----------------------------extraction---------------------------------------

def sliceTerm(term, p1, p2):
	# Removes item from between 2 punctuation marks
	cdef int idx
	cdef int ind
	idx = term.find(p1)
	if p1 == p2:
		ind = term.rfind(p2)
	else:
		ind = term.find(p2)
	if idx < ind:
		# Drop item in parentheses/quotes
		if ind == len(term)-1:
			term = term[:idx]
		elif idx == 0:
			term = term[ind+1:]
		else:
			term = term[:idx] + term[ind+1:]
	else:
		# Remove puntuation
		term = term.replace(p1, "")
		term = term.replace(p2, "")
	return term

def checkName(query):
	# Check query format
	cdef int idx
	cdef int ind
	cdef str term = query.lower()
	# Replace multiple spaces
	term = sub(r" +", " ", term)
	if "?" in term or "not " in term or "unknown" in term:
		# Skip uncertain entries
		return term, "uncertainEntry"
	if term[-2:] == " x" or term[-4:] == " mix" or "mix " in term or "hybrid " in term or " hybrid" in term: 
		return term, "hybrid"
	if "(" in term or ")" in term:
		term = sliceTerm(term, "(", ")")
	if "/" in term:
		# Subset from longer side of slash
		idx = term.find("/")
		if idx <= len(term)/2:
			term = term[idx+1:]
		else:
			term = term [:idx]
	if '"' in term:
		term = sliceTerm(term, '"', '"')
	if "&" in term:
		# Replace ampersand and add spaces if needed
		idx = term.find("&")
		if 0 < idx < len(term)-1:
			if term[idx+1] != " ":
				# Check second space first so index remains accurate
				term = term[:idx+1] + " " + term[idx+1:]
			if term[idx-1] != " ":
				term = term[:idx] + " " + term[idx:]
			term = term.replace("&", "and")
		else:
			# Remove leading/trailing ampersand
			term = term.replace("&", "")
	if "#" in term:
		idx = term.find("#")
		if idx < len(term)-1:
			# Drop symbol and any numbers
			if idx < len(term)/2:
				ind = term[idx:].find(" ") + idx
				term = term[ind+1:]
			else:
				ind = term.rfind(" ")
				term = term[:ind]			
	return term, ""

def checkForNum(query):
	# Determines if a query is primarily letters
	cdef int count = 0
	cdef int c
	cdef str j
	cdef str i
	cdef list splt
	cdef str term = query
	if " " in term:
		# Attempt to remove numbers
		splt = term.split()
		term = ""
		for j in splt:
			for i in digits:
				c = j.count(i)
			if c < len(j)/4.0:
				term += j + " "
		term = term.strip()
	# Calculate total number content after trimming
	for i in digits:
		count += term.count(i)	
	if term and count < len(query)/4.0:
		return term, ""
	else:
		return "", "numberContent"

def filterNames(outfile, misses, t, query, reas=""):
	# Filters query and attmepts to correct formatting
	cdef str head
	cdef str term = ""
	if t and not reas:
		if len(query) >= 3:
			# Filter names by type
			if t == "common":
				term, reas = checkForNum(query)
				if not reas:
					term, reas = checkName(term)
					if not reas and len(term) < 3:
						reas = "formatting"
			elif t == "scientific":
				term, reas = checkForNum(query)
				if not reas:
					for i in punctuation:
						if i != "." and i in term:
							# Pass if query contains punctuation other than period
							reas = "punctuation"
							term = ""
							break
		else:
			reas = "tooShort"
	if not reas:
		# Fix caps
		head = term[0].upper()
		term = head + term[1:].lower()
		writeResults(outfile, ("{},{},{}\n").format(query, term.strip(), t))
		return 1
	else:
		writeResults(misses, ("{},{}\n").format(query, reas))
		return 0

def assignNames(outfile, misses, query):
	# Generates classifier and sorts names
	cdef str i
	cdef str t = ""
	cdef int x
	cdef int p = 0
	cdef int r = 0
	classifier = getSequenceClassifier()
	print("\n\tFiltering species names...")
	for i in query:
		if len(i) >= 3:
			t = classifier.classify(nameFeatures(i))
			if "," in i:
				# Replace commas in original name to preserve formatting
				i = i.replace(",", " ")
			x = filterNames(outfile, misses, t, i)
		else:
			# Forward to filterNames to record misses
			x = filterNames(outfile, misses, t, i, "tooShort")
		if x == 1:
			p += 1
		else:
			r += 1
	print(("\tSuccessfully formatted {} entries.\n\t{} entries failed formatting.").format(p, r))

def sortNames(outfile, misses, common, scientific, query):
	# Determines if names must be sorted before filtering
	cdef str i
	cdef str t = ""
	cdef int x
	cdef int p = 0
	cdef int r = 0
	# Assign name type or get classifier
	if common == True:
		t = "common"
	elif scientific == True:
		t = "scientific"
	if t:
		print("\n\tFiltering species names...")
		for i in query:
			if "," in i:
				# Replace commas in original name to preserve formatting
				i = i.replace(",", " ")
			x = filterNames(outfile, misses, t, i)
			if x == 1:
				p += 1
			else:
				r += 1
		print(("\tSuccessfully formatted {} entries.\n\t{} entries failed formatting.").format(p, r))
	else:
		assignNames(outfile, misses, query)
