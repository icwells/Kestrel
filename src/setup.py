'''This script will cythonize the Kestrel package.'''

from distutils.core import setup
from Cython.Build import cythonize

KT = "kestrelTools.pyx"
SP = "scrapePages.pyx"
TS = "taxaSearch.pyx"

# Print blank lines to split output
print(("\n\tComipiling {}...\n").format(KT))
setup(ext_modules=cythonize(KT))
print(("\n\tComipiling {}...\n").format(SP))
setup(ext_modules=cythonize(SP))
print(("\n\tComipiling {}...\n").format(TS))
setup(ext_modules=cythonize(TS))
print()
