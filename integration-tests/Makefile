LOGSATSH_VERISON := 7.3.1
PWD              := $(shell pwd)

logstash:
	docker run --rm -p 9600:9600 \
		-v $(PWD)/pipeline:/usr/share/logstash/pipeline/ \
		-v $(PWD)/test.log:/tmp/test.log \
		-e XPACK_MONITORING_ENABLED=false \
		-e QUEUE_TYPE=persisted  \
		docker.elastic.co/logstash/logstash:$(LOGSATSH_VERISON)