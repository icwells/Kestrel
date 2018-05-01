'''This script will cythonize the Kestrel package.'''

from distutils.core import setup
from Cython.Build import cythonize

KT = "kestrelTools.pyx"
SP = "scrapePages.pyx"
SS = "seleniumSearch.pyx"
TS = "taxaSearch.pyx"

# Print blank lines to split output
for i in [KT, SP, SS, TS]:
	print(("\n\tComipiling {}...\n").format(i))
	setup(ext_modules=cythonize(i))
print()
