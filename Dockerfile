FROM golang:alpine AS build
WORKDIR /go/src/chatApp
COPY . .
COPY .env ./
RUN go build -o /go/bin/chatApp cmd/main.go

FROM scratch
COPY --from=build /go/bin/chatApp /go/bin/chatApp
ENTRYPOINT ["/go/bin/chatApp"]
