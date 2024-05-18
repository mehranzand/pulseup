TAG := "$$(git describe --abbrev=0 --tags)"

.PHONY: docker-build
docker-build:
	docker build --build-arg TAG=$(TAG) -t mehranzand/pulseup:$(TAG) -t mehranzand/pulseup:latest .

.PHONY: dev
dev:
	cd web && yarn dev & air && fg

.PHONY: web-dev
web-dev:
	cd web && yarn dev

.PHONY: clean
clean:
	rm -rf web/dist
	go clean -i