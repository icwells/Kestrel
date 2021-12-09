'''Defines TensorFlow model for common/scientifc name classifier'''

from argparse import ArgumentParser
from datetime import datetime
from getpass import getpass
import matplotlib.pyplot as plt
import mysql.connector
import numpy as np
from random import shuffle
import tensorflow as tf
import tensorflow_hub as hub
from tensorflow.keras.preprocessing.text import Tokenizer
from tensorflow.keras.preprocessing.sequence import pad_sequences

class Classifier():

	def __init__(self, args):
		self.database = "kestrelTaxonomy"
		self.db = None
		self.epochs = 20
		self.host = "localhost"
		self.hub = "https://tfhub.dev/google/nnlm-en-dim50/2"
		self.labels_test = []
		self.labels_train = []
		self.model = None
		self.outfile = "nlpModel"
		self.test = []
		self.train = []
		self.training_size = 10000
		self.username = args.u
		self.__connect__()
		self.__getDataSets__()

	def __connect__(self):
		# Connects to taxonomy database
		print()
		if not self.username:
			self.username = input("\tEnter MySQL username: ")
		# Get password
		password = getpass(prompt = "\tEnter MySQL password: ")
		try:
			# Connect to database
			self.db = mysql.connector.connect(
				host = self.host,
				database = self.database,
				user = self.username,
				password = password
			)
		except:
			print("\n\tIncorrect username or password. Exiting.")
			quit()

	def __getList__(self, column, table, name):
		# Extracts target column from table
		ret = []
		cursor = self.db.cursor()
		sql = ("SELECT DISTINCT({}) FROM {};").format(column, table)
		cursor.execute(sql)
		for i in cursor.fetchall():
			ret.append([name, i[0].strip()])
		return ret

	def __getDataSets__(self):
		# Extracts common and scientific lists from database
		train = []
		labels = []
		print("\n\tReading SQL tables...")
		names = self.__getList__("Name", "Common", 0)
		names.extend(self.__getList__("Species", "Taxonomy", 1))
		shuffle(names)
		# Get training and testing sets
		for i in names:
			# Split labels and terms after shuffling
			labels.append(i[0])
			train.append(i[1])
		self.labels_train = np.array(labels[:self.training_size])
		self.labels_test = np.array(labels[self.training_size:])
		self.train = np.array(train[:self.training_size])
		self.test = np.array(train[self.training_size:])

	def __plot__(self, history, metric):
		# Plots results
		plt.plot(history.history[metric])
		plt.plot(history.history['val_'+metric])
		plt.xlabel("Epochs")
		plt.ylabel(metric)
		plt.legend([metric, 'val_'+metric])
		plt.savefig("{}.svg".format(metric), format="svg")
  
	def trainModel(self):
		# Trains species name classifier
		print("\tTraining model...")
		hub_layer = hub.KerasLayer(self.hub, input_shape=[], dtype=tf.string, trainable=True)
		self.model = tf.keras.Sequential([
			hub_layer,
			tf.keras.layers.Dense(16, activation='relu'),
			tf.keras.layers.Dense(1, activation="sigmoid")
		])
		self.model.compile(loss='binary_crossentropy', optimizer='adam', metrics=['accuracy'])
		print(self.model.summary())
		history = self.model.fit(self.train, self.labels_train, 
				epochs = self.epochs, 
				batch_size = 512, 
				validation_data = (self.test, self.labels_test), 
				verbose = 2
		)
		print(self.model.evaluate(self.test, self.labels_test))
		self.__plot__(history, "accuracy")
		self.__plot__(history, "loss")

	def save(self):
		# Stores model in outfile
		tf.saved_model.save(self.model, self.outfile)

def main():
	start = datetime.now()
	parser = ArgumentParser("")
	parser.add_argument("-u", help = "MySQL username.")
	c = Classifier(parser.parse_args())
	c.trainModel()
	c.save()
	print(("\tTotal runtime: {}\n").format(datetime.now() - start))

if __name__ == "__main__":
	main()
