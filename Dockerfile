FROM golang:1.18 AS build

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN CGO_ENABLED=0 go build -o /usr/src/app/bin/app ./cmd/

FROM alpine:latest

WORKDIR /usr/src/app

COPY --from=build /usr/src/app/bin/app .
COPY --from=build /usr/src/app/configs/config.env .
COPY --from=build /usr/src/app/internal/templates/response_template.html .
COPY --from=build /usr/src/app/wait-for-it.sh .

RUN apk add --no-cache bash

EXPOSE 8001