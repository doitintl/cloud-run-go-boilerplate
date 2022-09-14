# syntax=docker/dockerfile:1

### Build Stage
FROM golang:1.18-bullseye as build

WORKDIR /go/src/service-name
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 go build -o /go/bin/service-name

### Run Stage
FROM gcr.io/distroless/static-debian11
COPY --from=build /go/bin/service-name /
COPY --from=build /go/src/service-name/config_data/ /config_data/
ENTRYPOINT ["/service-name"]

