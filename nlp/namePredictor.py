'''Predicts whether input names are common or scientific.'''

from argparse import ArgumentParser
from datetime import datetime
import numpy as np
import os
import tensorflow as tf

class Predictor():

	def __init__(self, args):
		self.infile = args.i
		self.model = tf.keras.models.load_model("nlpModel")
		self.names = []
		self.outfile = args.o
		self.__getNames__()
		self.__predict__()
		#self.__write__()

	def __getNames__(self):
		# Reads in single column of input names
		with open(self.infile, "r") as f:
			for line in f:
				self.names.append(line.strip())
		self.names = np.array(self.names)

	def __predict__(self):
		# Predicts whether name is common/scientific
		results = self.model.predict(self.names)
		print(results)

	def __write__(self):
		# Writes output to file
		with open(self.outfile, "w") as out:
			for i in self.names:
				out.write(",".join(i) + "\n")				

def main():
	start = datetime.now()
	parser = ArgumentParser("Predicts whether input names are common or scientific.")
	parser.add_argument("-i", help = "Path to input file.")
	parser.add_argument("-o", help = "Path to output file.")
	Predictor(parser.parse_args())
	print(("\tTotal runtime: {}\n").format(datetime.now() - start))

if __name__ == "__main__":
	main()
