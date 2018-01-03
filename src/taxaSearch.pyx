'''These functions will .'''

import nltk
from urllib import request, error
from kestrelTools import nameFeatures
from scrapePages import *

def getmatches(x, y):
	# Determines if taxonomies are equal
	if not x and not y:
		return False, None
	elif len(x) == len(y):
		for i in range(len(x)):
			if 

def searchTerms(outfile, misses, keys, common, query):
	# Serches for each item in terms list
	cdef str gu
	cdef str wu = "NA"
	cdef str match = ""
	cdef str term = query
	if len(query) > 1 and "," not in query:
		if common == True:
			gu = "NA"
			while len(term.split()) >= 1:
				# Search EOL and NCBI
				e, eu = searchEOL(term, keys[EOL])
				n, nu = searchNCBI(term, keys[NCBI])
				match = getmatches([e, n], [eu, nu])
				if not match:
					# Search Wikipedia to resolve mismatch
					w, wu = searchSource(WIKI, term)
					match = getmatches([e, n, w], [eu, nu, wu])
				if not match and len(term.split()) > 1:
					# Remove first word and try again
					term = term[term.find(" ")+1:]
				else:
					break
		elif common == False:
			# Search GBIF
			g, gu = searchSource(GBIF, term)
			e, eu = searchEOL(term, keys[EOL])
			n, nu = searchNCBI(term, keys[NCBI])
			match = getmatches([e, n, g], [eu, nu, gu])
			if not match:
				# Search Wikipedia
				w, wu = searchSource(term)
				match = getmatches([e, n, g], [eu, nu, gu])
	if match:
		# Write taxonomy
		with open(outfile, "a") as output:
			# Input name, taxonomy, source urls
			output.write(("{},{}\n").format(query, match)
	else:
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
