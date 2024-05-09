.PHONY: start startapp startdb createdb migrateup migratedown restartdb dropdb runcmd sqlc test get imagebuild imagerebuild tour server

#############################################################################################################################################################################
#																																											#
#	The following commands are general purpose.																																#
#																																											#
#############################################################################################################################################################################


all: start

start:
	docker-compose up

#############################################################################################################################################################################
#																																											#
#	The following commands are used to manage the database.																													#
#																																											#
#############################################################################################################################################################################

startdb:
	docker-compose up db \
	&& make createdb \
	&& make migrateup

createdb:
	docker-compose exec db createdb -U exampleuser simple_bank

dropdb:
	docker-compose exec db dropdb -U exampleuser simple_bank

migrateup:
	docker-compose run --rm app migrate -path ./db/migrations -database "postgresql://exampleuser:test1234@db:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	docker-compose run --rm app migrate -path ./db/migrations -database "postgresql://exampleuser:test1234@db:5432/simple_bank?sslmode=disable" -verbose down

restartdb: dropdb createdb migrateup

dbcmd:
	docker-compose exec db psql -U exampleuser -d simple_bank -c "$(cmd)"

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
	docker-compose run --rm app sh -c "mockgen -package mockdb -destination db/mock/$(dest).go github.com/LeonDavidZipp/go_simple_bank/db/sqlc $(iname)"

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
