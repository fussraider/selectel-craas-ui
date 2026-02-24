FROM alpine:latest as certs
RUN apk --no-cache add ca-certificates
RUN adduser -D -H -u 10001 appuser

FROM scratch
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=certs /etc/passwd /etc/passwd
COPY --from=certs /etc/group /etc/group
COPY server /server
USER appuser
EXPOSE 8080
ENTRYPOINT ["/server"]
