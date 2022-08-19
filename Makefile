run:
	go run ./cmd/go-commentof/ ./testdata/fixture
.PHONY: run

update-output:
	go run ./cmd/go-commentof/ ./testdata/fixture > ./testdata/output.json
	go run ./cmd/go-commentof/ -all ./testdata/fixture > ./testdata/output-all.json
.PHONY: update-output

view:
	go doc -all ./testdata/fixture
.PHONY: view