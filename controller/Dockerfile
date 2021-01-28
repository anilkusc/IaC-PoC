FROM golang:1.15 as BUILD
WORKDIR /src
COPY go.sum go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /bin/app .

FROM alpine:3.8
RUN mkdir -p /app/certs
WORKDIR /app
COPY --from=BUILD /bin/app .
RUN chmod +x app
ENTRYPOINT ./app