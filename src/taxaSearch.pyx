'''These functions will upload and extract data from the ASUviral database.'''

import nltk
import json
import codecs
from bs4 import BeautifulSoup
from collections import OrderedDict
from urllib import request, error
from kestrelTools import nameFeatures

cdef str EOL= "http://eol.org/api/"
cdef str SEARCH = "search/1.0."
cdef str PAGES = "pages/1.0."
cdef str HIER = "hierarchy_entries/1.0."
cdef str FORMAT = "xml"

cdef str GBIF = "http://api.gbif.org/v1/"
cdef str WIKI = "https://en.wikipedia.org/wiki/"
TAXONOMY = OrderedDict([("Kingdom",""),("Phylum",""),("Class",""),("Order",""),("Family",""),("Genus",""),("Species","")])

def scrapeWiki(soup):
	# Extract taxonomy data from Wikipedia entry
	cdef list ret = []
	cdef str k
	t = OrderedDict(TAXONOMY)
	for i in soup.find_all("tr"):
		# Iterate through template fields
		if i.td:
			k = ""
			for j in i.find_all("td"):
				if j.string:
					# Identify taxonomy values
					key = j.string.replace(":", "").strip()
					if key in t.keys():
						k = key
					if k:
						# Get taxonomy from next field
						if k == "Species":
							if j.span and j.span.i:
								if j.span.i.b and j.span.i.b.string:
									t[k] = j.span.i.b.string
									k = ""
									break
						elif j.a and j.a.string:
							t[k] = j.a.string
							k = ""
							break
	if t["Genus"]:
		for k in t.keys():
			# Convert to list if genus is present
			if t[k] and t[k] != " ":
				ret.append(t[k])
			else:
				ret.append("NA")
	return ret

def scrapeEOL(soup):
	# Scrapes EOL page for taxonomy
	cdef str b
	cdef list block
	cdef str line
	cdef str i
	cdef str name = ""
	cdef str rank = ""
	cdef int next = 0
	cdef list ret = []
	t = OrderedDict(TAXONOMY)
	b = "".join(soup.prettify())
	block = b.split("\n")
	for line in block:
		line = line.strip()
		if next == 1:
			# Get rank and name for own line
			if name:
				rank = line
			else:
				name = line
			next = 0
		if "<dwc:scientificname>" in line.lower():
			if "</" in line:
				# Get name from within field
				name = line[line.find(">")+1:line.rfind("<")]
			else:
				next = 1
		elif "<dwc:taxonrank>" in line.lower():
			if "</" in line:
				# Get rank from witihin field
				rank = line[line.find(">")+1:line.rfind("<")]
			else:
				next = 1
		if rank and name:
			# Add to dict if a linnaean classification and clear for next level
			rank = rank[0].upper() + rank[1:].lower()
			if rank in t.keys():
				if name.count(" ") > 1:
					n = name.split()
					name = n[0] + " " + n[1]
				t[rank] = name
			rank = ""
			name = ""
	if t["Genus"]:
		for i in t.keys():
			if t[i] and t[i] != " ":
				# Correct improper formatting before appending
				if i != "Species" and " " in t[i]:
					t[i] = t[i][:t[i].find(" ")]
				if "," in t[i]:
					t[i] = t[i].replace(",", "")
				ret.append(t[i])
			else:
				ret.append("NA")
	return ret

def scrapeGBIF(js):
	# Scrapes taxonomy from JSON output from GBIF
	cdef str i
	cdef list ret = []
	t = OrderedDict(TAXONOMY)
	if js["results"]:
		j = js["results"][0]
		for i in t.keys():
			if i.lower() in j.keys():
				t[i] = j[i.lower()]
	if t["Genus"]:
		for i in t.keys():
			if t[i]:
				ret.append(t[i])
			else:
				ret.append("NA")
	return ret

#-----------------------------------------------------------------------------

def getPage(source, term, keys=None):
	# Returns soup instance
	cdef str url
	# Format api query
	if source == EOL:
		url = ("{}{}{}?id={}&vetted=1&key={}").format(source, HIER, FORMAT, term, keys[EOL])
	else:
		# Check caps
		t = term[0].upper()
		term = t + term[1:].lower()
		if source == GBIF:
			if " " in term:
				# Percent encode spaces
				term = term.replace(" ", "%20")
			url = ("{}species?name={}").format(source, term)
		elif source == WIKI:
			if " " in term:
				# Replace spaces with underscores
				term = term.replace(" ", "_")
			url = source + term
	try:
		result = request.urlopen(url)
		return result, url
	except:
		return None, url

def getHID(query, key):
	# Gets hierarchy entry id from EOL
	cdef str url
	cdef str block
	url = ("{}{}{}?id={}&vetted=1&key={}").format(EOL, PAGES, FORMAT, query, key)
	try:
		result = request.urlopen(url)
		soup = BeautifulSoup(result, "lxml")
		for i in soup.find_all("taxon"):
			# Manually loop through (colon causes syntax error)
			block = "".join(i.prettify())
			try:
				start = block.find("<dwc:taxonid>") + len("<dwc:taxonid>") + 1
				end = block.find("</dwc:taxonid>")
				return block[start:end].strip()
			except IndexError:
				return None
	except:
		return None

def getTID(query, key):
	# Gets taxon ID from EOL
	cdef str url
	# Percent encode spaces
	query = query.replace(" ", "%20")
	url = ("{}{}{}?id={}&&vetted=1&key={}").format(EOL, SEARCH, FORMAT, query, key)
	try:
		result = request.urlopen(url)
		soup = BeautifulSoup(result, "lxml")
		for i in soup.find_all("entry"):
			if i.id and i.id.string:
				# Get first hit (no way to resolve multiples)
				return i.id.string
	except:
		return None

def searchSource(source, query, keys=None):
	# Searches given database for match
	cdef list t = []
	if source == EOL:
		# Get page id for species and extract page
		taxonid = getTID(query, keys[EOL])
		if taxonid:
			hierid = getHID(taxonid, keys[EOL])
			if hierid:
				result, url = getPage(source, hierid, keys)
				if result:
					# Remove api key
					url = url[:url.find("&")]
					t = scrapeEOL(BeautifulSoup(result, "lxml"))
	elif source == GBIF:
		result, url = getPage(source, query)
		if result:
			# Convert http bytes objet to string before loading into json
			reader = codecs.getreader("utf-8")
			t = scrapeGBIF(json.load(reader(result)))
	elif source == WIKI:
		result, url = getPage(source, query)
		if result:
			t = scrapeWiki(BeautifulSoup(result, "lxml"))
	if t and t.count("NA") <= 1:
		# Save query, taxonomy and source url
		return ("{},{}\n").format(",".join(t), url)
	else:
		return None

def searchTerms(outfile, misses, keys, common, query):
	# Serches for each item in terms list
	cdef str match = ""
	cdef str term = query
	if len(query) > 1 and "," not in query:
		if common == True:
			while len(term.split()) >= 1:
				# Search EOL
				match = searchSource(EOL, term,  keys)
				if match:
					break
				# Search Wikipedia
				match = searchSource(WIKI, term)
				if not match and len(term.split()) > 1:
					# Remove first word and try again
					term = term[term.find(" ")+1:]
				else:
					break
		elif common == False:
			# Search GBIF
			match = searchSource(GBIF, term)
			if not match:
				# Search EOL
				match = searchSource(EOL, term,  keys)
				if not match:
					# Search Wikipedia
					match = searchSource(WIKI, term)
	if match:
		# Write taxonomy
		with open(outfile, "a") as output:
			output.write(query + "," + match)
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
