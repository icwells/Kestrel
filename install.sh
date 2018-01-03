##############################################################################
#	Installs Go and Cython packages for Kestrel
#
#		Requires:	Cython
##############################################################################

KT="kestrelTools"
SP="scrapePages"
TS="taxaSearch"

cd src/
python setup.py build_ext --inplace
rm -r build/
rm *.c
cd ../

mv src/$KT.*.so bin/$KT.so
mv src/$SP.*.so bin/$SP.so
mv src/$TS.*.so bin/$TS.so
