FROM golang:1.21 as base

WORKDIR /app
COPY . /app
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build ./src/cmd/api/main.go

FROM alpine as prod
WORKDIR /app
COPY --from=base /app/main ./

EXPOSE 8080 8081
CMD [ "./main" ]
