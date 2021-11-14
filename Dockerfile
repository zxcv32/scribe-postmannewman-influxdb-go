FROM golang:1.17 AS build
WORKDIR /scribe
COPY go.mod /scribe
COPY go.sum /scribe
RUN go mod download
COPY *.go /scribe
RUN CGO_ENABLED=0 GOOS=linux go build -installsuffix cgo -tags=nomsgpack -a -o app .

FROM alpine:3
LABEL NAME="scribe-influxdb-go"
LABEL version="0.0.1"
WORKDIR /scribe
ENV GIN_MODE=release
EXPOSE 8090
COPY .env ./
COPY --from=build /scribe/app ./
# Please specify at least INFLUXDB_ORG and INFLUXDB_TOKEN at runtime
CMD ["./app"]
