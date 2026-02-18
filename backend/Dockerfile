FROM alpine:latest as certs
RUN apk --no-cache add ca-certificates

FROM scratch
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY server /server
EXPOSE 8080
ENTRYPOINT ["/server"]
