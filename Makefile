DOCKER ?= docker
SOLR_9_IMAGE ?= solr:9.7
SOLR_8_IMAGE ?= solr:9.7
SOLR_CLOUD_NAME ?= solrcloud
SOLR_NAME ?= solr

.PHONY: unit-test
unit-test:
	go test -v -cover

.PHONY: integration-test
integration-test:
	go test -tags integration -v -cover

.PHONY: solr-9
solr-9: rm-solr
	$(DOCKER) run -d -p 8983:8983 --name $(SOLR_NAME) $(SOLR_9_IMAGE) solr -f
	$(DOCKER) cp fixtures/security.json $(SOLR_NAME):/var/solr/data/security.json
	$(DOCKER) exec -t $(SOLR_NAME) bash -c 'sleep 10; wait-for-solr.sh --max-attempts 10 --wait-seconds 10'
	$(DOCKER) exec -t $(SOLR_NAME) bash -c 'SOLR_AUTH_TYPE="basic" SOLR_AUTHENTICATION_OPTS="-Dbasicauth=solr:SolrRocks" solr create -c searchengines'

.PHONY: solrcloud-9
solrcloud-9: rm-solrcloud	
	$(DOCKER) run -d -p 8984:8983 --name $(SOLR_CLOUD_NAME) $(SOLR_9_IMAGE) solr -c -f
	$(DOCKER) cp fixtures/security.json $(SOLR_CLOUD_NAME):/tmp/security.json
	$(DOCKER) exec -t $(SOLR_CLOUD_NAME) bash -c 'solr zk cp file:/tmp/security.json zk:/security.json -z localhost:9983'
	$(DOCKER) exec -t $(SOLR_CLOUD_NAME) bash -c 'sleep 10; wait-for-solr.sh --max-attempts 10 --wait-seconds 10'

.PHONY: solr-8
solr-8: rm-solr
	$(DOCKER) run -d -p 8983:8983 --name $(SOLR_NAME) $(SOLR_8_IMAGE) solr -f
	$(DOCKER) cp fixtures/security.json $(SOLR_NAME):/var/solr/data/security.json
	$(DOCKER) exec -t $(SOLR_NAME) bash -c 'sleep 10; wait-for-solr.sh --max-attempts 10 --wait-seconds 10'
	$(DOCKER) exec -t $(SOLR_NAME) bash -c 'SOLR_AUTH_TYPE="basic" SOLR_AUTHENTICATION_OPTS="-Dbasicauth=solr:SolrRocks" solr create -c searchengines'

.PHONY: solrcloud-8
solrcloud-8: rm-solrcloud	
	$(DOCKER) run -d -p 8984:8983 --name $(SOLR_CLOUD_NAME) $(SOLR_8_IMAGE) solr -c -f
	$(DOCKER) cp fixtures/security.json $(SOLR_CLOUD_NAME):/tmp/security.json
	$(DOCKER) exec -t $(SOLR_CLOUD_NAME) bash -c 'solr zk cp file:/tmp/security.json zk:/security.json -z localhost:9983'
	$(DOCKER) exec -t $(SOLR_CLOUD_NAME) bash -c 'sleep 10; wait-for-solr.sh --max-attempts 10 --wait-seconds 10'

.PHONY: rm-solrcloud
rm-solrcloud: 
	$(DOCKER) rm -f $(SOLR_CLOUD_NAME) || true

.PHONY: rm-solr
rm-solr:
	$(DOCKER) rm -f $(SOLR_NAME) || true

.PHONY: cleanup
cleanup: rm-solr rm-solrcloud