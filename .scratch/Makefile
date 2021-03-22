test:
	./test_script.sh

test-report: test
	sed -i "s/_$$(pwd|sed 's@/@\\/@g')/./g" coverage.tmp
	go tool cover -html=coverage.tmp
