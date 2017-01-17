ll: compile docker push

compile:
	CGO_ENABLED=0 go build -o metrics

docker:
	sudo docker build --no-cache=true --force-rm=true -t sjezewski/coral-metrics .

push:
	sudo docker push sjezewski/coral-metrics

