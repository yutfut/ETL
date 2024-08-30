.PHONY: dc_up
dc_up:
	docker-compose -f ./deployment/docker-compose.yml up

.PHONY: dc_down
dc_down:
	docker-compose -f ./deployment/docker-compose.yml down
