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

for I in $KT $SP $SS $TS; do
	mv src/$I*so bin/$I*so
done

echo ""
echo "Done"
echo ""
