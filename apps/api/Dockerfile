FROM golang:1.24-alpine as builder
RUN apk add build-base
WORKDIR /build
COPY ./go.* ./
RUN go mod download
COPY ./apps/api ./apps/api
RUN CGO_ENABLED=1 go build -o ./main ./apps/api/main.go


FROM alpine
WORKDIR /app
COPY --from=builder /build/main ./main
EXPOSE 1323
CMD ["./main"]
