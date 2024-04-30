.PHONY: docker-build
docker-build:
	docker compose up --build --force-recreate

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