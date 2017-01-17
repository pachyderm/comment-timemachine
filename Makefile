one-time-setup: docker loop

docker:
	cd pipelines/metrics && make

loop:
	yes | pachctl delete-all
	pachctl create-repo stream
	ls -1 pipelines/input | while read line; do \
		pachctl start-commit stream master; \
		pachctl delete-file stream master loremipsum.json; \
		pachctl put-file stream master /loremipsum.json -f pipelines/input/$$line; \
		pachctl finish-commit stream master; \
		sleep 1; \
	done;
	pachctl create-pipeline -f pipeline.json

