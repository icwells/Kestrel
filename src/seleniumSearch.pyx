'''This script will call selenium to Google search potential taxonomies'''

import re
from string import ascii_lowercase
from collections import OrderedDict
from urllib import request, error
from bs4 import BeautifulSoup
from selenium import webdriver
from selenium.webdriver.common.keys import Keys
from scrapePages import *

WIKI = "https://en.wikipedia.org/wiki/"
IUCN = "http://www.iucnredlist.org"
ITIS = "https://www.itis.gov/"
TAXONOMY = OrderedDict([("Kingdom","NA"),("Phylum","NA"),("Class","NA"),("Order","NA"),("Family","NA"),("Genus","NA"),("Species","NA"),("url","NA")])

def splitName(s):
	# Removes extra info from species name line
	cdef int end = 0
	for idx,i in enumerate(s[1:]):
		if i not in ascii_lowercase and i != " ":
			end = idx + 1
			break
	return s[:end]

def scrapeITIS(soup, url):
	# Parses results from ITIS page
	cdef str k = ""
	t = OrderedDict(TAXONOMY)
	t["url"] = url
	for tr in soup.find_all("tr"):
		for td in tr.find_all("td"):
			if k:
				if k == "Species" and td.string:
					# Subset binomial from line
					t[k] = splitName(td.string)
				elif td.a and td.a.string:
					# Get rank name
					t[k] = td.a.string
				k = ""
				break
			elif td.string:
				# Remove formatting and whitespace
				s = td.string.strip()
				if s in t.keys():
					# Get taxonomy rank
					k = s
	return checkTaxa(t)

def scrapeIUCN(soup, url):
	# Scrapes taxonomy from IUCN webpage
	cdef list k = []
	cdef list n = []
	cdef int sci = 0
	t = OrderedDict(TAXONOMY)
	t["url"] = url
	for tr in soup.find_all("tr"):
		for th in tr.find_all("th"):
			if th.string:
				if th.string in t.keys():
					# Get table headers (taxonomy ranks)
					k.append(th.string)
		for td in tr.find_all("td"):
			if td.strong and td.strong.string:
				if td.strong.string == "Scientific Name:":
					sci = 1
			elif td.span and td.span.span and sci == 1:
				if td.span.span.string:
					# Get species and genus names
					t["Species"] = td.span.span.string
					t["Genus"] = t["Species"].split()[0]
					break
			elif k and td.string:
				if len(n) < len(k):
					# Store taxa in order
					n.append(td.string)
	for i in range(len(n)):
		# Store remaining taxa
		t[k[i]] = n[i]
	return checkTaxa(t)

def openURL(url):
	# Opens target url
	try:
		return request.urlopen(url)
	except:
		return None

def parseURLS(urls):
	# Attempts to find taxonomy from given urls
	cdef str i
	t = None
	for i in urls:
		if "#" in i:
			# Remove article subheader
			i = i[:i.find("#")]
		# Request page
		result = openURL(i)
		if result:
			# Only proceed if there is a response
			if WIKI in i:
				t = scrapeWiki(BeautifulSoup(result, "lxml"),i)
			elif ITIS in i:
				t = scrapeITIS(BeautifulSoup(result, "lxml"), i)
			if IUCN in i:
				t = scrapeIUCN(BeautifulSoup(result, "lxml"), i)
		if t:
			break
	# Return empty if no match is found
	return t

def getURLS(soup):
	# Extracts usable urls from Google result
	cdef str x
	cdef list urls = []
	links = soup.findAll("a")
	for link in soup.find_all("a"):
		try:
			# Seperate urls
			u = re.split(":(?=http)",link["href"].replace("/url?q=",""))
			for i in u:
				# Filter out unusable links
				for x in [Wiki, IUCN, ITIS]:
					if x in i:
						urls.append(i)
		except KeyError:
			pass
	return urls

def getSearchResult(browser, term):
	# Searches Google for term
	browser.get("http://www.google.com")
	# Find the search box
	elem = browser.find_element_by_name("q")
	elem.send_keys(term + " taxonomy" + Keys.RETURN)
	soup = BeautifulSoup(browser.page_source, "lxml")
	browser.quit()
	return soup

def getBrowser(firefox):
	# Returns selenium browser instance
	if firefox:
		return webdriver.Firefox()
	else:
		return webdriver.Chrome()
