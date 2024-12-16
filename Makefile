DOCKER ?= docker
SOLR_IMAGE ?= solr:9.7
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
	$(DOCKER) -t $(SOLR_NAME) bash -c 'exec chmod -R 755 /var/solr/data'
	$(DOCKER) -t $(SOLR_NAME) bash -c 'exec chown -R solr:solr /var/solr/data'
	$(DOCKER) cp fixtures/security.json $(SOLR_NAME):/var/solr/data/security.json
	$(DOCKER) exec -t $(SOLR_NAME) bash -c 'sleep 5; wait-for-solr.sh --max-attempts 10 --wait-seconds 5'
	$(DOCKER) exec -t $(SOLR_NAME) bash -c 'SOLR_AUTH_TYPE="basic" SOLR_AUTHENTICATION_OPTS="-Dbasicauth=solr:SolrRocks" solr create -c searchengines'

.PHONY: solrcloud
solrcloud: rm-solrcloud	
	$(DOCKER) run -d -p 8984:8983 --name $(SOLR_CLOUD_NAME) $(SOLR_IMAGE) solr -c -f
	$(DOCKER) cp fixtures/security.json $(SOLR_CLOUD_NAME):/tmp/security.json
	$(DOCKER) exec -t $(SOLR_CLOUD_NAME) bash -c 'solr zk cp file:/tmp/security.json zk:/security.json -z localhost:9983'
	$(DOCKER) exec -t $(SOLR_CLOUD_NAME) bash -c 'sleep 5; wait-for-solr.sh --max-attempts 10 --wait-seconds 5'

.PHONY: rm-solrcloud
rm-solrcloud: 
	$(DOCKER) rm -f $(SOLR_CLOUD_NAME) || true

.PHONY: rm-solr
rm-solr:
	$(DOCKER) rm -f $(SOLR_NAME) || true

.PHONY: cleanup
cleanup: rm-solr rm-solrcloud