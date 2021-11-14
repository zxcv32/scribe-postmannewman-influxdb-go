FROM golang:1.17 AS build
WORKDIR /scriber
COPY go.mod /scriber
COPY go.sum /scriber
RUN go mod download
COPY *.go /scriber
RUN CGO_ENABLED=0 GOOS=linux go build -installsuffix cgo -tags=nomsgpack -a -o app .

FROM alpine:3
LABEL NAME="scriber-influxdb-go"
LABEL version="0.0.1"
WORKDIR /scriber
ENV GIN_MODE=debug
EXPOSE 8090
COPY .env ./
COPY --from=build /scriber/app ./
# Please specify at least INFLUXDB_ORG and INFLUXDB_TOKEN at runtime
CMD ["./app"]
