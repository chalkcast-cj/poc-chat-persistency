down:
	- docker-compose down --remove-orphans

up: down
	docker-compose up
