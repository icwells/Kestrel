'''These functions will .'''

import nltk
from string import punctuation, digits
from urllib import request, error
from kestrelTools import nameFeatures
from scrapePages import *

cdef str EOL= "http://eol.org/api/"
cdef str NCBI = "https://eutils.ncbi.nlm.nih.gov/entrez/eutils/"
cdef str GBIF = "http://api.gbif.org/v1/"
cdef str WIKI = "https://en.wikipedia.org/wiki/"

def formatMatch(d):
	# Returns taxonomy as formatted string
	cdef int mn = 8
	cdef str i
	cdef int c
	cdef list m
	for i in d.keys():
		# Get most complete match
		c = list(d[i].values()).count("NA")
		if c < mn:
			mn = c
			hit = d[i]
	m = list(hit.values())[:-1]
	# Sort urls from remaining matches
	for i in [EOL,NCBI,GBIF,WIKI]:
		if i in d.keys():
			m.append(d[i]["url"])
		else:
			m.append("NA")
	return ",".join(m)

def findMatches(d, keys, col):
	# Returns key of matched terms
	cdef int l = len(keys)
	cdef int i = 0
	idx = set()
	while i < l:
		if d[keys[i]][col] == d[keys[i+1]][col]:
			# Save unique indices of matches
			idx.add(keys[i])
			idx.add(keys[i+1])
		i += 1
	return list(idx)

def getmatches(d, last=0):
	# Returns string of best match and associated urls
	cdef int match = 0
	cdef int n
	cdef list keys = list(d.keys())
	cdef int length = len(keys)
	cdef list col = list(d[keys[0]].keys())
	if length > 1:
	# Sort through order, family, genus, and species backwards
		for i in col[3:-1:-1]:
			idx = findMatches(d, keys, i)
			if len(idx) == length:
				# Accept perfect match
				match = 1
				break
			elif length > 2 and len(idx)/length >= 0.5:
				# Accept majority match
				match = 1
				break
	if match == 0 and last == 0:
		# Quit if no match found
		return ""
	else:
		return formatMatch(d)

def writeResults(outfile, line):
	# Writes match to outfile and no match to misses
	with open(outfile, "a") as output:
		output.write(line)

#-----------------------------------------------------------------------------

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

def sourceDict(l):
	# Converts list of dicts to dict of dicts
	d = {}
	for i in l:
		# Remove empty values
		if type(i) == dict and i == {}:
			pass
		elif i != None:
			for u in [EOL,NCBI,GBIF,WIKI]:
				if u in i["url"]:
					d[u] = i
	return d

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
			res = sourceDict([e, n])
			if res:
				match = getmatches(res)
			if not match:
				# Search Wikipedia to resolve mismatch
				w = searchWiki(term)
				res = sourceDict([e, n, w])
				if term.count(" ") == 0:
					last = 1
				if res:
					match = getmatches(res, last)
			if not match and last == 0:
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
		res = sourceDict([g, e, n])
		if res:
			match = getmatches(res)
		if not match:
			# Search Wikipedia
			w = searchWiki(query)
			res = sourceDict([g, e, n, w])
			if res:
				match = getmatches(res, 1)
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
