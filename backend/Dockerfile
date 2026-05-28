FROM golang:1.23-alpine AS build
WORKDIR /app
COPY main.go .
RUN go build -o server main.go

FROM alpine:3.20
COPY --from=build /app/server /server
EXPOSE 8081
CMD ["/server"]
