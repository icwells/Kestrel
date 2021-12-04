'''Defines TensorFlow model for common/scientifc name classifier'''

from argparse import ArgumentParser
from datetime import datetime
from getpass import getpass
import matplotlib.pyplot as plt
import MySQLdb
import os
import tensorflow as tf
import tensorflow_hub as hub
import tensorflow_datasets as tfds
import unixpath

class Classifier():

	def __init__(self, args):
		self.common = []
		self.db = None
		self.model = None
		#self.outdir = args.o
		self.scientific = []
		self.username = args.u
		self.__connect__()
		self.__getDataSets__()

	def __connect__(self):
		# Connects to taxonomy database
		if not self.username:
			self.username = input("\tEnter MySQL username: ")
		# Get password
		password = getpass(prompt = "\tEnter MySQL password: ")
		try:
			# Connect to database
			self.db = MySQLdb.connect("localhost", self.username, password, "kestrelTaxonomy")
		except:
			print("\n\tIncorrect username or password. Exiting.")
			quit()

	def __getList__(self, column, table):
		# Extracts target column from table
		cursor = self.db.cursor()
		sql = ("SELECT DISTINCT({}) FROM {};").format(column, table)
		cursor.execute(common)
		self.common = cursor.fetchall()

	def __getDataSets__(self):
		# Extracts common and scientific lists from database
		self.common = self.__getList("Name", "Common")
		self.scientific = self.__getList("Species", "Taxonomy")


	#def train(self):
		# Trains species name classifier
		

def main():
	start = datetime.now()
	parser = ArgumentParser("")
	#parser.add_argument("-o", help = "Path to output dir for model.")
	parser.add_argument("-u", help = "MySQL username.")
	c = Classifier(parser.parse_args())
	#c.train()
	#c.write()
	print(("\tTotal runtime: {}\n").format(datetime.now() - start))

if __name__ == "__main__":
	main()
