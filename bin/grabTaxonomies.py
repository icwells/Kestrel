'''This script will assign taxonomies from a manually curated dataset to entries with matching common/scientific names.'''

from argparse import ArgumentParser
import os
from datetime import datetime

def writeMissed(fails, misses):
	# Write missed to file
	print("\tWriting unmatched queries to file...")
	with open(fails, "w") as out:
		for i in misses:
			out.write(i)

def matchTaxa(infile, outfile, taxa, head):
	# Writes matched taxonomies to passes and returns misses
	first = True
	count = 0
	misses = []
	print("\tMatching input queries to taxonomies...")
	with open(outfile, "w") as out:
		with open(infile, "r") as f:
			for line in f:
				if first == False:
					spli = line.split(",")
					# Get formatted search term
					n = spli[1].strip().lower()
					if n in taxa.keys():
						# [Query, term, taxonomy...]
						t = spli[:2]
						t.extend(taxa[n])
						out.write((",").join(t))
						count += 1
					else:
						misses.append(line)
				else:
					# Write header
					out.write(head)
					# Store header for misses
					misses.append(line)
					first = False
	return misses, count

def getTaxa(infile):
	# Returns header and dict of taxonomy with scientific and common names as keys
	first = True
	taxa = {}
	print(("\n\tReading input taxonomy from {}...").format(infile))
	with open(infile, "r") as f:
		for line in f:
			if first == False:
				t = line.split(",")
				n = t[1].lower().strip()
				s = t[8].lower().strip()
				t = t[2:]
				# Add one entry for scientific name
				taxa[s] = t
				if n != s:
					# Add a common name entry if present
					taxa[n] = t
			else:
				head = line
				first = False
	return head, taxa

def checkArgs(args):
	# Makes sure args are present
	if not args.i:
		print("\n\t[Error] Please specify an input file. Exiting.\n")
		quit()
	if not args.r:
		print("\n\t[Error] Please specify a reference taxonomy file. Exiting.\n")
		quit()
	if not args.r:
		print("\n\t[Error] Please specify an output directory. Exiting.\n")
		quit()
	for i in [args.i, args.r]:
		if not os.path.isfile(i):
			print(("\n\t[Error] Cannot find {}. Exiting.\n").format(i))
	# Format outdir
	if args.o[-1] != "/":
		args.o += "/"
	if not os.path.isdir(args.o):
		os.mkdir(args.o)
	# Get outfile names
	name = os.path.split(args.i)[1]
	name = name[:name.find(".")]
	return args.o + name + ".withTaxa.csv", args.o + name + ".unmatched.csv",

def main():
	starttime = datetime.now()
	parser = ArgumentParser("This script will assign taxonomies from a manually \
curated dataset to entries with matching common/scientific names.")
	parser.add_argument("-i", help = "Path to input file (output of kestrel.py --extract).")
	parser.add_argument("-r", 
help = "Path to reference file of manually currated taxonomies (output of previous kestrel search).")
	parser.add_argument("-o", help = "Path to output directory.")
	args = parser.parse_args()
	passes, fails = checkArgs(args)
	head, taxa = getTaxa(args.r)
	misses, l = matchTaxa(args.i, passes, taxa, head)
	writeMissed(fails, misses)
	print(("\n\tFound matches for {} of {} queries.").format(l, len(misses)+l))
	print(("\tFinished. Runtime: {}\n").format(datetime.now()-starttime))

if __name__ == "__main__":
	main()
