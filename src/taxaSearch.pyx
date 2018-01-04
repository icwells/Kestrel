'''These functions will .'''

import nltk
from string import punctuation, digits
from collections import Counter
from urllib import request, error
from kestrelTools import nameFeatures
from scrapePages import *

cdef str EOL= "http://eol.org/api/"
cdef str NCBI = "https://eutils.ncbi.nlm.nih.gov/entrez/eutils/"
cdef str GBIF = "http://api.gbif.org/v1/"
cdef str WIKI = "https://en.wikipedia.org/wiki/"

def mostCommon(l):
	# Returns most common name and number of occurances
	try:
		data = Counter(l)
		return data.most_common(1)[0]
	except TypeError:
		return None

def removeEmpty(l):
	# Returns list withoout empty entries
	for i in l:
		if not i or i == None:
			# Drop item if empty
			del i
	if len(l) >= 1:
		return l
	else:
		return None

def formatMatch(l):
	# Returns taxonomy as formatted string
	cdef str g = "NA"
	cdef str e = "NA"
	cdef str n = "NA"
	cdef str w = "NA"
	cdef str m = ""
	cdef int mn = 8
	cdef int hit = 0
	l = removeEmpty(l)
	if len(l) < 1 or l == None:
		return m
	elif len(l) > 1:
		# Pick most complete match
		for idx,i in enumerate(l):
			if type(i) == list:
				if i.count("NA") < mn:
					mn = i.count("NA")
					hit = idx
	try:
		m = ",".join(l[hit][:-1])
	except TypeError, IndexError:
		return m
	# Sort urls from remaining matches
	for i in l:
		if i:
			if EOL in i[-1]:
				e = i[-1]
			elif NCBI in i[-1]:
				n = i[-1]
			elif GBIF in i[-1]:
				g = i[-1]
			elif WIKI in i[-1]:
				w = i[-1]
	return ("{},{},{},{},{}").format(m, e, n, g, w)

def getmatches(l, last=0):
	# Returns string of best match and associated urls
	cdef int match = 0
	cdef int length
	cdef int n
	l = removeEmpty(l)
	length = len(l)
	if l == None or length < 1:
		return ""
	elif length > 1:
		# Transpose and reverse taxonomies; examine species and genus columns
		t = list(map(list, zip(l[:-2:-1])))
		for i in t:
			name = mostCommon(i)
			if name:
				n = int(name[1])
				if n == length:
					# Break if all taxonomies match
					match = 1
					break
				elif n >= 2 and length-2 <= n < length:
					# Break if there are at least two out of three/four matches
					for j in range(len(i)):
						if i[j] != name[0]:
							# Remove discordant taxonomy
							del l[j]
					match = 1
					break
	if match == 0 and last == 0:
		# Quit if no match found
		return ""
	else:
		return formatMatch(l)

def writeResults(outfile, line):
	# Writes match to outfile and no match to misses
	with open(outfile, "a") as output:
		output.write(line)

def checkName(query):
	# Check query format
	cdef str term = query
	if "?" in term or "not " in term.lower():
		# Skip uncertain entries
		return ""
	if term[-2:] == " X" or term[-4:] == " mix": 
		return ""
	if "'" in term:
		# Percent encode apsotrophes
		term = term.replace("'", "%27")
	if "(" in term or ")" in term:
		term = term.replace("(", "")
		term = term.replace(")", "")
	if "/" in term:
		# Subset from longer side of slash
		idx = term.find("/")
		if idx <= len(term)/2:
			term = term[idx+1:]
		else:
			term = term [:idx]
	if '"' in term:
		term = term.replace('"', '')
	if len(term) > 1:
		# Check caps
		t = term[0].upper()
		term = t + term[1:].lower()
	return term

def checkForNum(query):
	# Determines if a query is primarily letters
	cdef int count = 0
	cdef int idx
	cdef str i
	cdef str term = query
	for i in digits:
		if i in term:
			count += 1
			idx = term.find(i)
			if idx < len(term)-1:
				if term[idx+1] == " ":
					# Drop trailing number
					term = term[:idx]
					term = term.strip()
				elif (idx - 1) >= 0 and term[idx-1] == " ":
					# Drop leading number
					term = term[idx+2:]
	if count >= len(query)/4:
		return ""
	else:
		return term

def searchCommon(outfile, misses, keys, query):
	# Serches for mathces for common names
	cdef str term = ""
	cdef str match = ""
	cdef str reas = ""
	cdef int np = 1
	cdef last = 0
	term = checkForNum(query)
	if not term:
		reas = "numberContent"
	else:
		term = checkName(term)
		if not term:
			reas = "formatting"
	if len(term) > 1:
		while len(term.split()) >= 1:
			# Search EOL and NCBI
			e = searchEOL(term, keys[EOL])
			n = searchNCBI(term, keys[NCBI])
			match = getmatches([e,n])
			if not match:
				# Search Wikipedia to resolve mismatch
				w = searchWiki(term)
				if term.count(" ") == 0:
					last = 1
				match = getmatches([e, n, w], last)
			if not match and len(term.split()) > 1:
				# Remove first word and try again
				term = term[term.find(" ")+1:]
			else:
				break
	if match:
		writeResults(outfile, ("{},{},{}\n").format(query, term, match))
	else:
		if not reas:
			reas = "noMatch"
		writeResults(misses, ("{},{}\n").format(query, reas))

def searchSci(outfile, misses, keys, query):
	# Serches for mathces for scientific names
	cdef str term
	cdef str match = ""
	cdef str reas = ""
	term = checkForNum(query)
	if not term:
		reas = "numberContent"
	else:
		for i in punctuation:
			if i != "." and i in term:
				# Pass if query contains punctuation other than period
				reas = "punctuation"
				break
		# Fix caps
		t = term[0].upper()
		term = t + term[1:].lower()
	if len(term) > 1 and not reas:
		# Search GBIF
		g = searchGBIF(term)
		e = searchEOL(term, keys[EOL])
		n = searchNCBI(term, keys[NCBI])
		match = getmatches([e, n, g])
		if not match:
			# Search Wikipedia
			w = searchWiki(query)
			match = getmatches([e, n, g], 1)
	if match:
		writeResults(outfile, ("{},{},{}\n").format(query, term, match))
	else:
		if not reas:
			reas = "noMatch"
		writeResults(misses, ("{},{}\n").format(query, reas))

def assignQuery(outfile, misses, keys, classifier, query):
	# Determines whether query is a scientific or common name
	cdef str t
	if len(query) >= 3:
		t = classifier.classify(nameFeatures(query))
		if t == "common":
			searchCommon(outfile, misses, keys, query)
		elif t == "scientific":
			searchSci(outfile, misses, keys, query)
	else:
		with open(misses, "a") as m:
			m.write(query + "\n")
