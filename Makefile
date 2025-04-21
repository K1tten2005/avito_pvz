PKG := ./...
COVERAGE_HTML=coverage.html
COVERPROFILE_TMP=coverprofile.tmp

.PHONY: test coverage cover-html clean

test:
	go test -v $(PKG)

cover:
	go test -json ./... -coverprofile coverprofile_.tmp -coverpkg=./... ; \
    grep -v -e 'mocks.go' -e 'mock.go' -e 'docs.go' -e '_easyjson.go' -e 'gen_sql.go' coverprofile_.tmp > coverprofile.tmp ; \
    rm coverprofile_.tmp ; \
	go tool cover -html ${COVERPROFILE_TMP} -o  $(COVERAGE_HTML); \
    go tool cover -func ${COVERPROFILE_TMP}

view-coverage:
	open $(COVERAGE_HTML)

integration-test:
	go test -v ./internal/pkg/pvz/repo -tags=integration

clean:
	rm -f $(COVERAGE_FILE) $(COVERAGE_HTML) ${COVERPROFILE_TMP} 
