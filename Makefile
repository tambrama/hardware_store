.PHONY: docker-up docker-down docker-rebuild docker-logs

docker-up:
	docker-compose -f docker-compose.yml up --build -d

docker-down:
	docker-compose -f docker-compose.yml down

docker-rebuild:
	docker-compose -f docker-compose.yml down -v && \
	docker-compose -f docker-compose.yml up --build -d

docker-logs:
	docker-compose -f docker-compose.yml logs -f backend
