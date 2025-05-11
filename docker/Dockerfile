FROM docker.io/golang:1.23.3

WORKDIR /app
COPY ./go.mod ./go.sum ./
RUN go mod download

COPY ./ .

WORKDIR /app/cmd/asbwig
RUN CGO_ENABLED=0 GOOS=linux go build

EXPOSE 8085

CMD ["/app/cmd/asbwig/asbwig"]