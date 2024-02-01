# CONTAINER FOR BUILDING BINARY
FROM golang:1.21 AS build

# INSTALL DEPENDENCIES
RUN go install github.com/gobuffalo/packr/v2/packr2@v2.8.3
COPY go.mod go.sum /src/
RUN cd /src && go mod download

# BUILD BINARY
COPY . /src
RUN cd /src/db && packr2
RUN cd /src/tools/signer && make build

# CONTAINER FOR RUNNING BINARY
FROM alpine:3.18.0
COPY --from=build /src/tools/signer/dist/x1-signer /app/x1-signer
COPY --from=build /src/tools/signer/config/signer.config.toml /app/example.config.toml
RUN apk update && apk add postgresql15-client
EXPOSE 7001
CMD ["/bin/sh", "-c", "/app/x1-signer"]
