# from root repo: "docker build -t go-mylife-energy -f docker/Dockerfile ."

FROM golang:1.18 as build

WORKDIR /go/src/app
COPY . .

WORKDIR /go/src/app/mylife-energy
RUN go mod download
# RUN go vet -v
# RUN go test -v

RUN CGO_ENABLED=0 go build -o /go/bin/web cmd/web/*.go

FROM gcr.io/distroless/static-debian11

COPY --from=build /go/bin/web /
COPY --from=build /go/src/app/mylife-energy/config.yaml /
CMD ["/web"]