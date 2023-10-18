FROM golang:1.21

WORKDIR /app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY . .
RUN go mod download && go mod verify

COPY *.go ./
RUN go build -o /resume-server

EXPOSE 8080

CMD ["/resume-server"]
