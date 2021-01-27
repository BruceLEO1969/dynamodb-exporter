FROM golang as builder

WORKDIR /go/src/github.com/BruceLEO1969/dynamodb-exporter
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o dynamodb-exporter .

FROM alpine

COPY --from=builder /go/src/github.com/BruceLEO1969/dynamodb-exporter/dynamodb-exporter /bin
RUN apk --update add --no-cache ca-certificates

EXPOSE 9436

CMD ["dynamodb-exporter"]
