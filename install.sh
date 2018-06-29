##############################################################################
#	Installs Go and Cython packages for Kestrel
#
#		Requires:	Cython
##############################################################################

KT="kestrelTools"
SP="scrapePages"
SS="seleniumSearch"
TS="taxaSearch"

cd src/
python setup.py build_ext --inplace
rm -r build/
rm *.c
cd ../

mv src/$KT.so bin/$KT.so
mv src/$SP.so bin/$SP.so
mv src/$SS.so bin/$SS.so
mv src/$TS.so bin/$TS.so

echo ""
echo "Done"
echo ""
