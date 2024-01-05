FROM golang:alpine AS build
WORKDIR /go/src/chatApp
COPY . .
COPY .env ./
RUN go build -o /go/bin/chatApp cmd/main.go

COPY ca-bundle.crt /etc/ssl/certs/ca-bundle.crt
COPY ca-bundle.trust.crt /etc/ssl/certs/ca-bundle.trust.crt 

FROM scratch
COPY --from=build /go/bin/chatApp /go/bin/chatApp
ENTRYPOINT ["/go/bin/chatApp"]
