# Kestrel Taxonomy Finder Version 0.3

## Copyright 2017 by Shawn Rupp

### Kestrel is a program for resolving species’ common names and synonyms with "official" scientific names and extracting taxonomies from internet databases.

## Dependencies
Python3  
Cython  
NLTK  
BeautifulSoup4 
Selenium 
Chrome or Firefox 

## Installation

### Cython
Kestrel utilizes Cython to compile python code into C and drastically improve performance. Cython can be installed from the pypi repository or via Miniconda (it is installed by default with the full Anaconda package).

	To install with Miniconda:
	conda install cython

### NLTK
Kestrel uses python’s Natural Language Processing Toolkit to differentiate between common and scientific names in its input. To install on any Debian-based Linux platform, enter the following into a terminal:

	conda install nltk
or
	sudo pip install -U nltk

Kestrel comes with it’s own training dataset, so you do not need to download any additional data from NLTK. 

### BeautifulSoup4
Kestrel also uses the BeautifulSoup module, and the lxml parser, to parse hmtl and xml pages.

	apt-get install python3-bs4
	apt-get install python-lxml

### Selenium
Kestrel uses Selenium to run automated Google searches. To install the Python repository, enter the following into a terminal:

	pip install -U selenium

Selenium requires browser-specific drivers. To install, select either the Chrome or Firefox driver on Seleium's pypi page (https://pypi.python.org/pypi/selenium) and follow the given instructions.

### Kestrel
Download the git repository, change into the directory, and build the Cython scripts.

	git clone https://github.com/icwells/Kestrel.git
	cd Kestrel/
	./install.sh

## For further documentation, see KestrelReadMe.pdf
