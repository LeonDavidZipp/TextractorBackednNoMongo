include .env

.PHONY: start startapp startdb createdb migrateup migratedown restartdb dropdb runcmd sqlc test get imagebuild imagerebuild tour server

#############################################################################################################################################################################
#																																											#
#	The following commands are general purpose.																																#
#																																											#
#############################################################################################################################################################################

all: start

start:
	docker-compose up

stop:
	docker-compose down

tidy:
	docker-compose run --rm backend go mod tidy

#############################################################################################################################################################################
#																																											#
#	The following commands are used to manage the user database.																											#
#																																											#
#############################################################################################################################################################################

startdb:
	docker-compose up userdb

# creates fresh database && tables in already running db container
createdb:
	docker-compose exec userdb createdb -U $(POSTGRES_USER) $(POSTGRES_DB_NAME)

dropdb:
	docker-compose exec userdb dropdb -U $(POSTGRES_USER) $(POSTGRES_DB_NAME)

resetdb: dropdb createdb migrateup

migrateup:
	docker-compose run --rm backend migrate -path ./db/migrations -database "$(POSTGRES_SOURCE)" -verbose up

migratedown:
	docker-compose run --rm backend migrate -path ./db/migrations -database "$(POSTGRES_SOURCE)" -verbose down

dbcmd:
	docker-compose exec userdb psql -U exampleuser -d $(POSTGRES_DB_NAME) -c "$(cmd)"


#############################################################################################################################################################################
#																																											#
#	The following commands are used to run (in) the application.																											#
#																																											#
#############################################################################################################################################################################

startapp:
	docker-compose up backend
	
sqlc:
	docker-compose run --rm backend sh -c "sqlc generate"

test:
	docker-compose run --rm backend sh -c "go test -v -cover ./..."

get:
	docker-compose run --rm backend go get -u $(pkg)

server:
	docker-compose run --rm backend sh -c "go run main.go"

appcmd:
	docker-compose run --rm backend sh -c "$(cmd)"

mockdb:
	docker-compose run --rm backend sh -c "mockgen -package mockdb -destination db/mock/$(dest).go github.com/LeonDavidZipp/Textractor/db/store $(iname)"


#############################################################################################################################################################################
#																																											#
#	The following commands are used to handle docker tasks																													#
#																																											#
#############################################################################################################################################################################

imagebuild:
	docker-compose build

imagerebuild:
	docker-compose build --no-cache

#############################################################################################################################################################################
#																																											#
#	A tour about the syntax of go																																			#
#																																											#
#############################################################################################################################################################################

tour:
	tour

