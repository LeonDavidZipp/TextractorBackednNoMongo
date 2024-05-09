########################################################################################
#                                                                                      #
#    Base                                                                              #
#                                                                                      #
########################################################################################

FROM golang:1.22.2-alpine3.19 as base

LABEL maintainer="lzipp"

WORKDIR /app
COPY ./db ./db
# COPY ./api ./api
COPY ./go.mod ./
COPY ./go.sum ./
RUN go mod download

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# CMD ["go", "run", "main.go"]

# EXPOSE 8080


########################################################################################
#                                                                                      #
#    Development                                                                       #
#                                                                                      #
########################################################################################

FROM base as dev

COPY ./sqlc.json ./

RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest \
    && go install go.uber.org/mock/mockgen@latest


########################################################################################
#                                                                                      #
#    Production                                                                        #
#                                                                                      #
########################################################################################

# FROM base as prod
