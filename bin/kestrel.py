'''This script contains functions for searching online databases for taxonomy information.'''

import argparse
from datetime import datetime, date
from functools import partial
from multiprocessing import Pool, cpu_count
from sys import stdout
from os import remove
from kestrelTools import *
from taxaSearch import assignQuery
from seleniumSearch import searchMisses

def checkArgs(args):
	# Makes sure proper arguments have been specified
	end = False
	if not args.i:
		print("\n\t[Error] Please provide an input file.")
		end = True
	if not args.o:
		if end != True:
			# Print newline to split from command prompt
			print()
		print("\t[Error] Please provide an output file.")
	if args.extract or args.merge:
		if args.c < 0:
			if end != True:
				print()
				end = True
			print("\t[Error] Please provide a column number.")			
	if end == True:
		print("\tExiting.\n")
		quit()

def version():
	print("\n\tKestrel v0.6 (8/13/2018) is a program for resolving common names and synonyms with \
scientific names and extracting taxonomies.")
	print("\n\tCopyright 2017 by Shawn Rupp.")
	print("\tThis program comes with ABSOLUTELY NO WARRANTY.\n\tThis is free \
software, and you are welcome to redistribute it under certain conditions.\n")
	quit()

def main():
	starttime = datetime.now()
	parser = argparse.ArgumentParser(description = "Kestrel will search online \
databases for taxonomy information.")
	parser.add_argument("-v", action = "store_true", 
help = "Prints version info and exits.")
	parser.add_argument("--extract", action = "store_true", default = False,
help = "Extracts and filters input names.")
	parser.add_argument("--common", action = "store_true", default = False,
help = "Indicates that input contains only common names (use with --extract).")
	parser.add_argument("--scientific", action = "store_true", default = False,
help = "Indicates that input contains only scientific names (use with --extract).")
	parser.add_argument("-c", type = int, default = -1, help = "Column containing \
species names (integer starting from 0; use with --extract)).")
	parser.add_argument("-i", help = "Path to input file.")
	parser.add_argument("-o", help = "Path to output csv file.")
	parser.add_argument("-t", default = 1, type = int,
help = "Number of threads for identifying taxa (default = 1).")
	parser.add_argument("--firefox", action = "store_true", default = False,
help = "Use Firefox browser (uses Chrome by default).")
	parser.add_argument("--merge", action = "store_true", default = False,
help = "Merges output taxonomy (given with -o) with original input file (-i). Column of species names \
in input file must be given with -c. Output will be written in same directory as taxonomy file.") 
	args = parser.parse_args()
	if args.v:
		version()
	else:
		checkArgs(args)
	if args.extract == True:
		print("\n\tExtracting and filtering species names...")
		done = checkOutput(args.o, "Query,SearchTerm,Type\n")
		misses = args.o[:args.o.rfind("/")+1] + "KestrelRejected.csv"
		missed = checkOutput(misses, "Query,Reason\n")
		query = speciesList(args.i, args.c, done.extend(missed))
		sortNames(args.o, misses, args.common, args.scientific, query)
	elif args.merge == True:
		name = os.path.split(args.i)[1]
		print(("\n\tMerging taxonomies from {} with {}...").format(os.path.split(args.o)[1], name))
		outfile = args.o[:args.o.rfind("/")+1] + name[:name.find(".")] + ".withTaxonomies.csv"
		taxa = getTaxa(args.o)
		mergeTaxonomy(args.i, outfile, args.c, taxa)
	else:
		match = 0
		miss = 0
		print("\n\tGenerating taxonomy output...")
		if args.t > cpu_count():
			args.t = cpu_count()
		# Get input data
		keys = apiKeys()
		misses = args.o[:args.o.rfind("/")+1] + "KestrelMisses.csv"
		header = "Query,SearchTerm,Kingdom,Phylum,Class,Order,Family,Genus,Species,\
EOL,NCBI,Wikipedia,IUCN,GBIF,ITIS\n"
		done = checkOutput(args.o, header)
		missed = checkOutput(misses, "Query,SearchTerm,Reason\n")
		# Store missed and done lengths
		donelen = len(done)
		done.extend(missed)
		missedlen = len(missed)
		# Read in query names
		query = termList(args.i, done)
		l = float(len(query)) + float(len(done))
		pool = Pool(processes = args.t)
		func = partial(assignQuery, args.o, misses, keys)
		# API search
		print(("\n\tIdentifying species with {} threads....").format(args.t))
		for i,x in enumerate(pool.imap_unordered(func, query)):
			stdout.write("\r\t{0:.1%} of query names have finished".format((i+donelen+missedlen)/l))
			if x > 0:
				match += x
			else:
				miss += abs(x)
		pool.close()
		pool.join()
		print(("\n\tFound matches for {} entries.").format(match + donelen))
		print(("\tNo match found for {} entries.").format(miss + missedlen))
		# Google search
		print("\n\tSearching for missed terms...")
		nomatch = args.o[:args.o.rfind("/")+1] + "KestrelNoMatch.csv"
		nm = checkOutput(nomatch, "Query,SearchTerm,Reason\n")
		newquery = termList(misses)
		hits, nohit = searchMisses(args.firefox, args.o, nomatch, newquery)
		if hits:
			# Delete temp misses file
			remove(misses)
		print(("\n\tTotal matches found: {}").format(match + donelen + hits))
		print(("\tTotal entries without matches: {}").format(nohit))
	print(("\tFinished. Runtime: {}\n").format(datetime.now()-starttime))

if __name__ == "__main__":
	main()
