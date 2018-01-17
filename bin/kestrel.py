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
		if end != True:
			# Print newline to split from command prompt
			print()
		print("\t[Error] Please provide an output file.")
	if args.extract and args.c < 0:
		if end != True:
			print()
		print("\t[Error] Please provide a column number.")			
	if end == True:
		print("\n\tExiting.\n")
		quit()

def version():
	print("\n\tKestrel v0.2 (~) is a program for resolving common names and synonyms with \
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
	parser.add_argument("--extract", action = "store_true",
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
	args = parser.parse_args()
	if args.v:
		version()
	else:
		checkArgs(args)
	if args.extract:
		print("\n\tExtracting and filtering species names...")
		done = checkOutput(args.o, "Query,SearchTerm,Type\n")
		misses = args.o[:args.o.rfind("/")+1] + "KestrelRejected.csv"
		missed = checkOutput(misses, "Query,Reason\n")
		query = speciesList(args.i, args.c, done)
		sortNames(args.o, misses, args.common, args.scientific, query)
	else:
		print("\n\tGenerating taxonomy output...")
		if args.t > cpu_count():
			args.t = cpu_count()
		keys = apiKeys()
		misses = args.o[:args.o.rfind("/")+1] + "KestrelMisses.csv"
		header = "Query,SearchTerm,Kingdom,Phylum,Class,Order,Family,Genus,Species,EOL,NCBI,GBIF,Wikipedia\n"
		done = checkOutput(args.o, header)
		missed = checkOutput(misses)
		done.extend(missed)
		# Read in query names
		query = speciesDict(args.i, done)
		l = float(len(query)) + float(len(done))
		pool = Pool(processes = args.t)
		func = partial(assignQuery, args.o, misses, keys)
		print(("\n\tIdentifying species with {} threads....").format(args.t))
		for i, _ in enumerate(pool.imap_unordered(func, query), 1):
			stdout.write("\r\t{0:.1%} of query names have finished".format(i/l))
		pool.close()
		pool.join()
	print(("\n\tFinished. Runtime: {}\n").format(datetime.now()-starttime))

if __name__ == "__main__":
	main()
