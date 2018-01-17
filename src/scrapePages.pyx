'''This script contains webpage specific functions for extracting taxonomies'''

import json
import codecs
from bs4 import BeautifulSoup
from collections import OrderedDict
from urllib import request, error

cdef str EOL= "http://eol.org/api/"
cdef str SEARCH = "search/1.0."
cdef str PAGES = "pages/1.0."
cdef str HIER = "hierarchy_entries/1.0."
cdef str FORMAT = "xml"

cdef str NCBI = "https://eutils.ncbi.nlm.nih.gov/entrez/eutils/"
cdef str GBIF = "http://api.gbif.org/v1/"
cdef str WIKI = "https://en.wikipedia.org/wiki/"
TAXONOMY = OrderedDict([("Kingdom","NA"),("Phylum","NA"),("Class","NA"),("Order","NA"),("Family","NA"),("Genus","NA"),("Species","NA"),("url","NA")])

def getPage(source, term, key=None):
	# Returns soup instance
	cdef str url
	# Format api query
	if source == EOL:
		url = ("{}{}{}?id={}&vetted=1&key={}").format(source, HIER, FORMAT, term, key)
	else:
		if source == GBIF:
			url = ("{}species?name={}").format(source, term)
		elif source == WIKI:
			if "%20" in term:
				# Replace spaces with underscores
				term = term.replace("%20", "_")
			url = source + term
	try:
		result = request.urlopen(url)
		return result, url
	except:
		return None, url

def checkTaxa(t):
	# Checks dict content before returning
	cdef list v
	v = list(t.values())
	if v.count("NA") <= 2 and t["Genus"] != "NA":
		# Convert to list if genus is present
		for i in t.keys():
			# Correct improper formatting before appending
			t[i] = str(t[i])
			if "," in t[i]:
				t[i] = t[i].replace(",", "")
			if i == "Species":
				if t[i] != "NA":
					# Enforce binomial name
					if " " not in t[i]:
						t[i] = ("{} {}").format(t["Genus"], t[i])
			elif i != "Species" and " " in t[i]:
				t[i] = t[i][:t[i].find(" ")]
		return t
	else:
		return {}

def scrapeWiki(soup, url):
	# Extract taxonomy data from Wikipedia entry
	cdef list ret = []
	cdef str k
	t = OrderedDict(TAXONOMY)
	t["url"] = url
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
									t[k] = str(j.span.i.b.string)
									k = ""
									break
						elif j.a and j.a.string:
							t[k] = j.a.string
							k = ""
							break
	if t["Species"] != "NA":
		# Remove genus abbreviation (full genus name added later)
		t["Species"] = t["Species"][3:]
	return checkTaxa(t)

def searchWiki(query):
	# Searches Wikipedia
	result, url = getPage(WIKI, query)
	if result:
		return scrapeWiki(BeautifulSoup(result, "lxml"), url)
	else:
		return None

def scrapeGBIF(js, url):
	# Scrapes taxonomy from JSON output from GBIF
	cdef str i
	cdef list ret = []
	t = OrderedDict(TAXONOMY)
	t["url"] = url
	if js["results"]:
		j = js["results"][0]
		for i in t.keys():
			if i.lower() in j.keys():
				t[i] = j[i.lower()]
	return checkTaxa(t)

def searchGBIF(query):
	# Searches GBIF for synonyms
	result, url = getPage(GBIF, query)
	if result:
		# Convert http bytes object to string before loading into json
		reader = codecs.getreader("utf-8")
		return scrapeGBIF(json.load(reader(result)), url)
	else:
		return None

#------------------------------NCBI-------------------------------------------

def efetch(query, key):
	# Fetches taxonomy for given id
	cdef str url
	cdef str k
	cdef list ret = []
	t = OrderedDict(TAXONOMY)
	url = ("{}efetch.fcgi?db=Taxonomy&id={}$retmode={}&key={}").format(NCBI, query, FORMAT, key)
	t["url"] = url
	try:
		result = request.urlopen(url)
	except:
		return None
	soup = BeautifulSoup(result, "lxml")
	# Extract Linnaean taxonomy from NCBI taonomy page
	for i in soup.find_all("taxon"):
		if i.scientificname and i.rank:
			if i.scientificname.string and i.rank.string:
				rank = i.rank.string
				# Convert to title capitalization
				rank = rank[0].upper() + rank[1:].lower()
				if rank in t.keys():
					t[rank] = i.scientificname.string
	return checkTaxa(t)

def esearch(source, query, key):
	# Searches for species ID
	query = query.replace(" ", "%20")
	url = ("{}esearch.fcgi?db=Taxonomy&term={}&key={}").format(source, query, key)
	try:
		result = request.urlopen(url)
		soup = BeautifulSoup(result, "lxml")
		for i in soup.find_all("id"):
			if i.string:
				return i.string
	except:
		return None

def espell(source, term, key):
	# Checks for corrected spelling
	term = term.replace(" ", "%20")
	url = ("{}espell.fcgi?db=Taxonomy&term={}&key={}").format(source, term, key)
	try:
		# Check term spelling
		result = request.urlopen(url)
		soup = BeautifulSoup(result, "lxml")
		for i in soup.find_all("correctedquery"):
			if i.string:
				return i.string
	except:
		# Return original term if no match is found
		return term

def searchNCBI(term, key):
	# Searches NCBI
	query = ""
	idx = ""
	if term:
		query = espell(NCBI, term, key)
	if query:
		idx = esearch(NCBI, query, key)
	if idx:
		return efetch(idx, key)
	else:
		return None

#------------------------------EOL--------------------------------------------

def scrapeEOL(soup, url):
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
	t["url"] = url
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
	return checkTaxa(t)

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

def searchEOL(term, key):
	# Searches EOL for taxonomy
	# Get page id for species and extract page
	taxonid = getTID(term, key)
	if taxonid:
		hierid = getHID(taxonid, key)
		if hierid:
			result, url = getPage(EOL, hierid, key)
			if result:
				# Remove api key
				url = url[:url.find("&")]
				return scrapeEOL(BeautifulSoup(result, "lxml"), url)
	else:
		return None
