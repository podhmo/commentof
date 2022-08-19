run:
	go run ./cmd/go-commentof/ ./testdata/fixture
.PHONY: run

update-output:
	go run ./cmd/go-commentof/ ./testdata/fixture > ./testdata/output.json
	go run ./cmd/go-commentof/ -all ./testdata/fixture > ./testdata/output-all.json
.PHONY: update-output

check-output:
	rm -f testdata/ng.*
	( grep unexported testdata/output.json && touch testdata/ng.unexported ) || :
	( grep _test.go testdata/output.json && touch testdata/ng.testfile ) || :
	test -z "`ls testdata/ng.* 2>/dev/null`"
.PHONY: check-output

view:
	go doc -all ./testdata/fixture
.PHONY: view