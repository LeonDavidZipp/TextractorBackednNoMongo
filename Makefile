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
	docker-compose run --rm app go mod tidy

#############################################################################################################################################################################
#																																											#
#	The following commands are used to manage the user database.																											#
#																																											#
#############################################################################################################################################################################

startdb:
	docker-compose up db

# creates fresh database && tables in already running db container
createdb:
	docker-compose exec db createdb -U $(POSTGRES_USER) $(POSTGRES_DB_NAME)

dropdb:
	docker-compose exec db dropdb -U $(POSTGRES_USER) $(POSTGRES_DB_NAME)

resetdb: dropdb createdb migrateup

migrateup:
	docker-compose run --rm app migrate -path ./db/migrations -database "$(POSTGRES_SOURCE)" -verbose up

migratedown:
	docker-compose run --rm app migrate -path ./db/migrations -database "$(POSTGRES_SOURCE)" -verbose down

dbcmd:
	docker-compose exec db psql -U exampleuser -d $(POSTGRES_DB_NAME) -c "$(cmd)"

startmongo:
	docker-compose up mongo-db

# createmongo

#############################################################################################################################################################################
#																																											#
#	The following commands are used to run (in) the application.																											#
#																																											#
#############################################################################################################################################################################

startapp:
	docker-compose up app
	
sqlc:
	docker-compose run --rm app sh -c "sqlc generate"

test:
	docker-compose run --rm app sh -c "go test -v -cover ./..."

get:
	docker-compose run --rm app go get -u $(pkg)

server:
	docker-compose run --rm app sh -c "go run main.go"

appcmd:
	docker-compose run --rm app sh -c "$(cmd)"

mockdb:
	docker-compose run --rm app sh -c "mockgen -package mockdb -destination db/mock/$(dest).go github.com/LeonDavidZipp/Textractor/db/store $(iname)"

startmdb:
	docker-compose up mongo-db

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

