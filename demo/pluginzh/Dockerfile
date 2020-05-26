FROM golang:1.14 as build

WORKDIR /workspace

COPY . /workspace

RUN CGO_ENABLED=0 go build -o bin/greeter pluginzh/greeter_impl.go

FROM alpine

ENV PLUGIN_MIN_PORT=7000 \
    PLUGIN_MAX_PORT=7000 \
    PLUGIN_HOST=true

COPY --from=build /workspace/bin/greeter /greeter

CMD ["/greeter"]

