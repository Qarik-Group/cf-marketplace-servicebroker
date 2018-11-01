FROM golang:1.11.1 as build
WORKDIR /go/src/github.com/starkandwayne/cf-marketplace-servicebroker/
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go install -v github.com/starkandwayne/cf-marketplace-servicebroker/cmd/cf-marketplace-servicebroker

FROM alpine:latest as final
RUN apk --no-cache add ca-certificates bash
WORKDIR /root/
EXPOSE 8080
ENV PORT 8080
COPY --from=build /go/bin/cf-marketplace-servicebroker .
CMD ["./cf-marketplace-servicebroker"]
