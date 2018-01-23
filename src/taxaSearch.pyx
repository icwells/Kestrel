'''These functions will .'''

from urllib import request, error
from kestrelTools import *
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

#-----------------------------------------------------------------------------

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

def searchCommon(outfile, misses, keys, query, term):
	# Serches for mathces for common names
	cdef str match = ""
	cdef int last = 0
	cdef int total = 0
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
			if len(term.split()) == 1:
				last = 1
			if res:
				match = getmatches(res, last)
		if not match and last == 0:
			# Remove first word and try again
			term = term[term.find(" ")+1:]
		else:
			break
	if match:
		# Replace percent formatting
		term = term.replace("%20", " ").replace("%27", "'")
		for i in query:
			total += 1
			writeResults(outfile, ("{},{},{}\n").format(i, term, match))
		# Return number of matched queries
		return total
	else:
		for i in query:
			total += 1
			writeResults(misses, ("{},noMatch\n").format(i))
		# Return negative to indicate failed queries
		return 0-total

def searchSci(outfile, misses, keys, query, term):
	# Serches for mathces for scientific names
	cdef str match = ""
	cdef str i
	cdef int total = 0
	# Search GBIF
	g = searchGBIF(term)
	e = searchEOL(term, keys[EOL])
	n = searchNCBI(term, keys[NCBI])
	res = sourceDict([g, e, n])
	if res:
		match = getmatches(res)
	if not match:
		# Search Wikipedia
		w = searchWiki(term)
		res = sourceDict([g, e, n, w])
		if res:
			match = getmatches(res, 1)
	if match:
		# Replace percent formatting
		term = term.replace("%20", " ").replace("%27", "'")
		for i in query:
			total += 1
			writeResults(outfile, ("{},{},{}\n").format(i, term, match))
		# Return number of matched queries
		return total
	else:
		for i in query:
			total += 1
			writeResults(misses, ("{},noMatch\n").format(i))
		# Return negative to indicate failed queries
		return 0-total

def assignQuery(outfile, misses, keys, query):
	# Determines whether query is a scientific or common name
	cdef int x
	cdef list q
	cdef str t = query[0]
	cdef str term = t
	if len(query) >= 3:
		# Isolate and type cast list of common names mapping to search term
		q = query[2:]
		if " " in term:
			# Percent encode spaces
			term = term.replace(" ", "%20")
		if "'" in term:
			# Percent encode apsotrophes
			term = term.replace("'", "%27")
		if query[1] == "common":
			x = searchCommon(outfile, misses, keys, q, term)
		elif query[1] == "scientific":
			x = searchSci(outfile, misses, keys, q, term)
	return x
