'''These functions will .'''

import nltk
from string import punctuation
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

def getmatches(l, last=1):
	# Returns string of best match and associated urls
	cdef str g = "NA"
	cdef str e = "NA"
	cdef str n = "NA"
	cdef str w = "NA"
	cdef int match = 0
	cdef str m = ""
	cdef int mn = 8
	cdef int hit
	cdef int length
	l = removeEmpty(l)
	length = len(l)
	if l == None or length < 1:
		return m
	# Transpose and reverse taxonomies
	t = list(map(list, zip(l[::-1])))
	for i in t:
		name = mostCommon(i)
		if name:
			if name[1] == length:
				# Break when all taxonomies match match (sorting backwards from most specific)
				match = 1
				break
			elif name[1] >= 2 and name[1] == length-1:
				# Break if there are at least two out of three matches
				for j in range(len(i)):
					if [j] != name[0]:
						# Remove discordant taxonomy
						del l[j]
				match = 1
				break
	if match == 0 and last == 0:
		# Quit if no match found
		return m
	else:
		l = removeEmpty(l)
		if len(l) < 1 or l == None:
			return m
		elif match == 1:
			# Join taxonomy 
			m = ",".join(l[0][:-1])
		elif match == 0 and last == 1:
			# Pick most complete match
			for idx,i in enumerate(l):
				if type(i) == list:
					if i.count("NA") < mn:
						mn = i.count("NA")
						hit = idx
			if hit < len(l):
				try:
					m = ",".join(l[hit][:-1])
				except TypeError, IndexError:
					return m
			else:
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

def searchTerms(outfile, misses, keys, common, query):
	# Serches for each item in terms list
	cdef str term = query
	cdef str match = ""
	cdef int np = 1
	for i in punctuation:
		if i != "-" and i != ".":
			if i in query:
				# Pass if query contains punctuation
				np = 0
				break
	if len(query) > 1 and np == 1:
		if common == True:
			last = 0
			while len(term.split()) >= 1:
				# Search EOL and NCBI
				e = searchEOL(term, keys[EOL])
				n = searchNCBI(term, keys[NCBI])
				match = getmatches([e,n])
				if not match:
					# Search Wikipedia to resolve mismatch
					w = searchWiki(term)
					if len(term.split()) == 1:
						last = 1
					match = getmatches([e, n, w], last)
				if not match and len(term.split()) > 1:
					# Remove first word and try again
					term = term[term.find(" ")+1:]
				else:
					break
		elif common == False:
			# Search GBIF
			g = searchGBIF(term)
			e = searchEOL(term, keys[EOL])
			n = searchNCBI(term, keys[NCBI])
			match = getmatches([e, n, g])
			if not match:
				# Search Wikipedia
				w = searchWiki(term)
				match = getmatches([e, n, g], True)
	if match:
		# Write taxonomy
		with open(outfile, "a") as output:
			# Input name, taxonomy, source urls
			output.write(("{},{}\n").format(query, match))
	elif query:
		# Write unmatched name
		with open(misses, "a") as m:
			m.write(query + "\n")

def assignQuery(outfile, misses, keys, classifier, query):
	# Determines whether query is a scientific or common name
	cdef str t
	if len(query) >= 3:
		t = classifier.classify(nameFeatures(query))
		if t == "common":
			searchTerms(outfile, misses, keys, True, query)
		elif t == "scientific":
			searchTerms(outfile, misses, keys, False, query)
	else:
		with open(misses, "a") as m:
			m.write(query + "\n")
