FROM golang:1.13.8 AS build-env

ADD . /dockerdev
WORKDIR /dockerdev

RUN go build -o /server

# Final Stage
FROM nginx:latest

EXPOSE 80

WORKDIR /

COPY --from=build-env /dockerdev/config.yaml /config.yaml
COPY --from=build-env /server /

CMD ["/server"]