'''Defines TensorFlow model for common/scientifc name classifier'''

from argparse import ArgumentParser
from datetime import datetime
from getpass import getpass
#from langdetect import detect
import matplotlib.pyplot as plt
import mysql.connector
import numpy as np
import os
from random import shuffle
import tensorflow as tf
import tensorflow_hub as hub
from tensorflow.keras.preprocessing.text import Tokenizer
from tensorflow.keras.preprocessing.sequence import pad_sequences
import unixpath

embedding_dim = 32
epochs = 20
max_len = 5
oov_tok = "<OOV>"
training_size = 10000
trunc_type = "post"
vocab_size = 10000
padding_type = "post"

class Classifier():

	def __init__(self, args):
		self.database = "kestrelTaxonomy"
		self.db = None
		self.host = "localhost"
		self.labels_test = []
		self.labels_train = []
		self.model = None
		self.names = []
		self.outfile = "utils/nlpModel"
		self.padded_test = []
		self.padded_train = []
		self.tokenizer = None
		self.username = args.u
		self.__connect__()
		self.__getDataSets__()

	def __write__(self):
		# Writes test data
		with open("test.csv", "w") as out:
			for idx, i in enumerate(self.padded_train):
				out.write("{},{}\n".format(self.labels_train[idx], i))
		quit()

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

	def __getTokenizer__(self):
		# Tokenizes and pads names
		print("\tInitializing tokenizer...")
		train = []
		labels = []
		for i in self.names:
			# Split labels and terms after shuffling
			labels.append(i[0])
			train.append(i[1])
		self.labels_train = np.array(labels[:training_size])
		self.labels_test = np.array(labels[training_size:])
		self.tokenizer = Tokenizer(oov_token = oov_tok)
		self.tokenizer.fit_on_texts(train)
		train = pad_sequences(self.tokenizer.texts_to_sequences(train), padding = padding_type, maxlen = max_len, truncating = trunc_type)
		self.padded_train = np.array(train[:training_size])
		self.padded_test = np.array(train[training_size:])

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
		print("\n\tReading SQL tables...")
		self.names = self.__getList__("Name", "Common", 0)
		self.names.extend(self.__getList__("Species", "Taxonomy", 1))
		shuffle(self.names)
		# Get training and testing sets
		self.__getTokenizer__()

	def __plot__(self, history, metric):
		# Plots results
		plt.plot(history.history[metric])
		plt.plot(history.history['val_'+metric])
		plt.xlabel("Epochs")
		plt.ylabel(metric)
		plt.legend([metric, 'val_'+metric])
		plt.show()
  
	def train(self):
		# Trains species name classifier
		print("\tTraining model...")
		model = "https://tfhub.dev/google/nnlm-en-dim50/2"
		hub_layer = hub.KerasLayer(model, input_shape=np.shape(self.padded_train[:3]), dtype=tf.int32, trainable=True)
		print(hub_layer(self.padded_train[:3]))
		quit()

		self.model = tf.keras.Sequential([
			hub_layer,
			#tf.keras.layers.Embedding(training_size + 1, embedding_dim, input_length = max_len),
			#tf.keras.layers.SimpleRNN(max_len),
			#tf.keras.layers.Bidirectional(tf.keras.layers.LSTM(64, return_sequences=True)),
			#tf.keras.layers.Bidirectional(tf.keras.layers.LSTM(32)),
			#tf.keras.layers.Lamda(detect) # add language detection
			#tf.keras.layers.Dense(max_len, activation='relu'),
			tf.keras.layers.Dense(1, activation="sigmoid")
		])
		self.model.compile(loss='binary_crossentropy', optimizer='adam', metrics=['accuracy'])
		print(self.model.summary())
		history = self.model.fit(self.padded_train, self.labels_train, 
				epochs=epochs, 
				batch_size=512, 
				validation_data=(self.padded_test, self.labels_test), 
				verbose=1
		)
		#results = model.evaluate(self.padded_test, self.labels_test)
		#print(history)
		#self.__plot__(history, "accuracy")
		#self.__plot__(history, "loss")

	def save(self):
		# Stores model in outfile
		tf.saved_model.save(self.model, self.outfile)

def main():
	start = datetime.now()
	parser = ArgumentParser("")
	parser.add_argument("-u", help = "MySQL username.")
	c = Classifier(parser.parse_args())
	c.train()
	#c.save()
	print(("\tTotal runtime: {}\n").format(datetime.now() - start))

if __name__ == "__main__":
	main()
