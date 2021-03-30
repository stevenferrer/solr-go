DOCKER ?= docker
SOLR_IMAGE ?= solr:8.8
SOLR_CLOUD_NAME ?= solrcloud
SOLR_NAME ?= solr

.PHONY: unit-test
unit-test:
	go test -v -cover

.PHONY: integration-test
integration-test:
	go test -tags integration -v -cover

.PHONY: solr
solr: rm-solr
	$(DOCKER) run -d -p 8983:8983 --name $(SOLR_NAME) $(SOLR_IMAGE) solr -f
	$(DOCKER) exec -it $(SOLR_NAME) bash -c 'sleep 5; wait-for-solr.sh --max-attempts 10 --wait-seconds 5'
	$(DOCKER) exec -it $(SOLR_NAME) solr create -c searchengines

.PHONY: solrcloud
solrcloud: rm-solrcloud	
	$(DOCKER) run -d -p 8984:8983 --name $(SOLR_CLOUD_NAME) $(SOLR_IMAGE) solr -c -f
	$(DOCKER) exec -it $(SOLR_CLOUD_NAME) bash -c 'sleep 5; wait-for-solr.sh --max-attempts 10 --wait-seconds 5'

.PHONY: rm-solrcloud
rm-solrcloud: 
	$(DOCKER) rm -f $(SOLR_CLOUD_NAME) || true

.PHONY: rm-solr
rm-solr:
	$(DOCKER) rm -f $(SOLR_NAME) || true

.PHONY: cleanup
cleanup: rm-solr rm-solrcloud