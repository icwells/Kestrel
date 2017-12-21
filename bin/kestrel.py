'''This script contains functions for searching Wikipedia for taxonomy information.'''

import argparse
from datetime import datetime, date
from functools import partial
from multiprocessing import Pool, cpu_count
from sys import stdout
from kestrelTools import *
from taxaSearch import *

def checkArgs(args):
	# Makes sure proper arguments have been specified
	end = False
	if not args.i:
		print("\n\t[Error] Please provide an input file.")
		end = True
	if not args.o:
		if end == True:
			print("\t[Error] Please provide an output file.")
		else:
			print("\n\t[Error] Please provide an output file.")
			end = True		
	if args.c < 0:
		if end == True:
			print("\t[Error] Please provide a column number")			
		else:
			print("\n\t[Error] Please provide a column number.")
			end = True
	if end == True:
		print("\n\tExiting.\n")
		quit()

def version():
	print("\n\tKestrel v0.1 (12/21/17) is a program for resolving common names and synonyms with \
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
	parser.add_argument("--common", action = "store_true", 
help = "Indicates that input contains only common names.")
	parser.add_argument("--scientific", action = "store_true", 
help = "Indicates that input contains only scientific names.")
	parser.add_argument("-i", help = "Path to input file.")
	parser.add_argument("-o", help = "Path to output file.")
	parser.add_argument("-c", type = int, default = -1,
help = "Column containing species names (integer starting from 0).")
	parser.add_argument("-t", default = 1, type = int,
help = "Number of threads for identifying taxa (default = 1).")
	args = parser.parse_args()
	if args.v:
		version()
	else:
		checkArgs(args)
	print("\n\tGenerating taxonomy output...")
	if args.t > cpu_count():
		args.t = cpu_count()
	keys = apiKeys()
	misses = args.o[:args.o.rfind("/")+1] + "KestrelMisses.txt"
	header = "Query,Kingdom,Phylum,Class,Order,Family,Genus,Species,URL\n"
	done = checkOutput(args.o, header)
	missed = checkOutput(misses, "Query")
	done.extend(missed)
	# Read in query and target names
	query = speciesList(args.i, args.c, done)
	d = float(len(done))
	l = float(len(query)) + d
	pool = Pool(processes = args.t)
	if args.common:
		func = partial(searchTerms, args.o, misses, keys, True)
	elif args.scientific:
		func = partial(searchTerms, args.o, misses, keys, False)
	else:
		classifier = getSequenceClassifier()
		func = partial(assignQuery, args.o, misses, keys, classifier)
	print(("\n\tIdentifying species with {} threads....").format(args.t))
	for i, _ in enumerate(pool.imap_unordered(func, query), 1):
		stdout.write("\r\t{0:.1%} of query names have finished".format(i/l))
	pool.close()
	pool.join()
	print(("\n\tFinished. Runtime: {}\n").format(datetime.now()-starttime))

if __name__ == "__main__":
	main()
