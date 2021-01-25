PODMAN ?= "podman"
SOLR ?="solr-go" 
COLLECTION ?= "gettingstarted"

.PHONY: solr
solr: stop-solr
	$(PODMAN) run -d -p 8983:8983 --name $(SOLR) solr:latest solr-precreate $(COLLECTION)

.PHONY: stop-solr
stop-solr:
	$(PODMAN) rm -f $(SOLR_INST) || true