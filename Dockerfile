FROM cgr.dev/chainguard/go:1.19 as build

WORKDIR /work
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o ./conncheck


FROM cgr.dev/chainguard/static:latest

ENV HOME /home/nonroot

COPY --from=build /work/conncheck /conncheck
ENTRYPOINT ["/conncheck"]
