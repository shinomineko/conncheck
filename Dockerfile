FROM cgr.dev/chainguard/go:1.20 as build

WORKDIR /work
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o ./conncheck

FROM cgr.dev/chainguard/static:latest

COPY --from=build /work/conncheck /conncheck
EXPOSE 8080
ENTRYPOINT ["/conncheck"]
