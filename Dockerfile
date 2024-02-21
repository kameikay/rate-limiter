FROM golang:1.21.5 as builder
WORKDIR /app
COPY . .
RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-w -s" -o app cmd/app/main.go

FROM scratch
COPY --from=builder /app/app .
CMD ["./app"]