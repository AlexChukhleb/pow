docker-build:
	set DOCKER_BUILDKIT=0
	docker build --no-cache -f scripts/docker/DockerfileServer -t pow-server .
	docker build --no-cache -f scripts/docker/DockerfileClient -t pow-client .
	
