# build stage
FROM golang:alpine AS build-env
ADD . /go/src/app
RUN cd /go/src/app && go build -o docker-entrypoint

# final stage
FROM alpine
COPY --from=build-env /go/src/app/docker-entrypoint /
ENTRYPOINT /docker-entrypoint
