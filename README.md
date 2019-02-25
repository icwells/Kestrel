# Kestrel Taxonomy Finder Version  

## Copyright 2019 by Shawn Rupp  

### Kestrel is a program for resolving species’ common names and synonyms with "official" scientific names and extracting taxonomies from internet databases.  
### Kestrel is meant to reduce the amount of manual searching required for a project, but its results may still require some manual curation.  

## Dependencies  
Go 1.11+  
Xvfb  
Chrome or Firefox    

## Installation  

### Xvfb  
Selenium requires Xvfb to run. To install on Linux:  

	sudo apt-get install xvfb  

### Kestrel  
Download the git repository, change into the directory, and install.  

	git clone https://github.com/icwells/Kestrel.git  
	cd Kestrel/  
	./install.sh  

#### Testing  
Once you have installed the program and its dependencies, you may with to run the test script:  

	./test.sh

If everything is properly configured, it should not throw any errors.  

## For further documentation, see KestrelReadMe.pdf
