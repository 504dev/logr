FROM node:20.11-alpine AS frontend

COPY ./ /opt/logr
WORKDIR /opt/logr/frontend
RUN yarn install && yarn build && yarn cache clean


FROM golang:1.13-alpine AS gobuild

COPY --from=frontend /opt/logr /opt/logr
WORKDIR /opt/logr
RUN go build -o logr-server ./cmd/server/main.go


# Start fresh from a smaller image
FROM alpine:3.9

RUN apk add --no-cache git
COPY --from=gobuild /opt/logr /opt/logr
WORKDIR /opt/logr
ENTRYPOINT ./logr-server --config="./config.yml"
EXPOSE 7776 7778
