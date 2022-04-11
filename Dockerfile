FROM golang:1.18-buster AS build

WORKDIR /app

COPY go.* ./
RUN go mod download
RUN go mod verify

FROM build AS rest
COPY . ./
RUN go build -v -o server ./cmd/rest

FROM debian:buster-slim as final-rest
RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY --from=rest /app/server /app/server

LABEL traefik.enable=true
LABEL traefik.http.routers.rider-rest.rule=PathPrefix(`/api/riders`)
LABEL traefik.http.routers.rider-rest.entrypoints=web
LABEL traefik.http.routers.rider-rest.middlewares='serviceheaders, traefik-forward-auth'
LABEL traefik.http.middlewares.serviceheaders.headers.accesscontrolalloworiginlist=*
LABEL traefik.http.middlewares.serviceheaders.headers.accessControlAllowMethods='GET, POST'
LABEL traefik.http.middlewares.serviceheaders.headers.accessControlAllowHeaders='authorization, content-type'

EXPOSE 1234

CMD ["/app/server"]