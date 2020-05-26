FROM golang:1.14 as build

WORKDIR /workspace

COPY . /workspace

RUN CGO_ENABLED=0 go build -o bin/basic main.go

FROM alpine

COPY --from=build /workspace/bin/basic /basic

CMD ["/basic"]

