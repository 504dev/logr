FROM node:20.11-alpine AS frontend

COPY ./ /opt/logr
WORKDIR /opt/logr/frontend
RUN yarn install && yarn build && yarn cache clean

###

FROM golang:1.13
COPY --from=frontend /opt/logr /opt/logr
WORKDIR /opt/logr
RUN go build -o logr-server ./cmd/server/main.go

ENTRYPOINT ["./logr-server"]
CMD ["--config=","./config.yml"]
