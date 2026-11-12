# Staging image: the auth server plus the generated site (copied in by scripts/deploy-staging.sh).
FROM golang:1.25-alpine AS build
WORKDIR /src
COPY go.mod main.go ./
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o /staging-server .

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /staging-server /staging-server
COPY site /site
ENTRYPOINT ["/staging-server"]
