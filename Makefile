PODMAN ?= "podman"
SOLR ?="solr-go" 

.PHONY: unit-test
unit-test:
	go test -v -cover

.PHONY: integration-test
integration-test:
	go test -tags integration -v -cover

.PHONY: start-solr
start-solr: stop-solr
	$(PODMAN) run -d -p 8983:8983 --name $(SOLR) solr:8.8 solr -c -f
	$(PODMAN) exec -it $(SOLR) bash -c 'sleep 5; wait-for-solr.sh --max-attempts 10 --wait-seconds 5'

.PHONY: stop-solr
stop-solr:
	$(PODMAN) rm -f $(SOLR) || true

