PODMAN ?= "podman"
SOLR ?="solr-go" 
COLLECTION ?= "searchengines"

.PHONY: unit-test
unit-test:
	go test -v -cover -race

.PHONY: integration-test
integration-test:
	go test -tags integration -v -cover -race

.PHONY: solr
solr: stop-solr
	$(PODMAN) run -d -p 8983:8983 --name $(SOLR) solr:latest solr-precreate $(COLLECTION)

.PHONY: stop-solr
stop-solr:
	$(PODMAN) rm -f $(SOLR) || true