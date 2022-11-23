run:
	docker-compose up -d --build

stop:
	docker-compose down

db-logs:
	docker-compose logs -f db

movies-app-logs:
	docker-compose logs -f movies-app

migrates-up:
	migrate \
		-source file://docker/migrates \
		-database postgres://postgres:goLANGninja@localhost:5432/movies?sslmode=disable \
		up

migrates-down:
	migrate \
		-source file://docker/migrates \
		-database postgres://postgres:goLANGninja@localhost:5432/movies?sslmode=disable \
		down