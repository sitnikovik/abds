.DEFAULT_GOAL = run

.PHONY: run
run:
	docker-compose up -d --build

.PHONY: stop
stop:
	docker-compose down
