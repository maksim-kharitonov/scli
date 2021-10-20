build()
{
	# Versioning variables
	MAJOR=$(git describe --abbrev=0 --always)
	VERSION=$(git describe --long | tr - .)
	BRANCH=${CI_COMMIT_REF_NAME:-$(git symbolic-ref -q --short HEAD)}

    git archive HEAD --prefix=scli/ --format=tar.gz -o scli-${MAJOR}.tar.gz
    mv ./scli-${MAJOR}.tar.gz ~/rpmbuild/SOURCES/
    rpmbuild -ba build/package/scli.spec --define "_version ${VERSION}" --define "_major ${MAJOR}" --define "_branch ${BRANCH}"
}

usage()
{
	echo "Usage: $0 [build|help]"
}

##### Main
case "$1" in
	'build' )
		build
		;;
	'help' )
		usage
		exit
		;;
	* )
		usage
		exit 1
esac
