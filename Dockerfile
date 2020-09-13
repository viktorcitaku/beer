# Compile stage
FROM golang:1.15.2-alpine AS build-env
ADD . /dockerdev
WORKDIR /dockerdev
RUN CGO_ENABLED=0 go build -o /server

# Final stage
FROM alpine
EXPOSE 8080
ENV STATIC_FILES /opt/web
WORKDIR /
COPY --from=build-env /server /opt
COPY ./web /opt/web
CMD ["/opt/server"]