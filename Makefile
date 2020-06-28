SOLR_INST ?="solr-go-test" 
COLLECTION ?= "gettingstarted"

.PHONY: solr
solr: stop-solr
	docker run -d -p 8983:8983 --name $(SOLR_INST) solr:latest solr-precreate $(COLLECTION)

.PHONY: stop-solr
stop-solr:
	docker rm -f $(SOLR_INST) || true