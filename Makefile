TAG := "$$(git describe --abbrev=0 --tags)"

.PHONY: docker-build
docker-build:
	docker build --build-arg TAG=$(TAG) -t mehranzand/pulseup:$(TAG) .

.PHONY: run-dev
dev:
	cd web && yarn dev & air && fg

.PHONY: web-dev
web-dev:
	cd web && yarn dev

.PHONY: clean
clean:
	rm -rf web/dist
	go clean -i