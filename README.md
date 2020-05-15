[![Build Status](https://travis-ci.com/icwells/Kestrel.svg?branch=master)](https://travis-ci.com/icwells/Kestrel)

# Kestrel Taxonomy Finder

## Copyright 2019 by Shawn Rupp  

### Kestrel is a program for downloading taxonomies from internet databases given species’ common or scientific names.  
### Kestrel is meant to reduce the amount of manual searching required for a project, but its results may still require some manual curation.  

## Dependencies  
[Go 1.11+](https://golang.org/dl/)  
Xvfb  
Chrome browser    

## Installation  

### Xvfb  
Selenium requires Xvfb to run. To install on Linux:  

	sudo apt-get install xvfb  

### Kestrel  
Download the git repository, change into the directory, and install (required Go packages will be installed).  

	git clone https://github.com/icwells/Kestrel.git  
	cd Kestrel/  
	./install.sh all  

#### Testing  
Once you have installed the program and its dependencies, you may wish to run the test script:  

	./test.sh

If everything is properly configured, it should not throw any errors.  

## For further documentation, see [KestrelReadMe.pdf](https://github.com/icwells/Kestrel/blob/master/KestrelReadMe.pdf)
