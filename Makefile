TALK_SESSION_SECRET=abbbcadabraeeeeeeeee

one-time-setup: docker input-data

initialize-submodules:
	git submodule init
	git submodule update --init
	cd talk && npm install
	cd proxy && npm install

docker:
	cd pipelines/metrics && make

input-data:
	yes | pachctl delete-all
	pachctl create-repo stream
	ls -1 pipelines/input | while read line; do \
		pachctl start-commit stream master; \
		pachctl delete-file stream master loremipsum.json; \
		pachctl put-file stream master /loremipsum.json -f pipelines/input/$$line; \
		pachctl finish-commit stream master; \
		sleep 1; \
	done;
	pachctl create-pipeline -f pipelines/pipeline.json

run: initialize-submodules
	# Check that pachyderm is running and we can connect
	which pachctl
	pachctl version
	# This will fail if redis/mongo not started
	(cd talk && echo $$PWD && TALK_SESSION_SECRET=bbbbaaaaabaaa npm start &)
	(cd proxy && node server.js &)

