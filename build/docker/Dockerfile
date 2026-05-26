# use specific versions of images
FROM openapitools/openapi-generator-cli:latest AS api

WORKDIR /local

COPY api api
COPY scripts scripts
COPY .openapi-generator-ignore .openapi-generator-ignore

RUN docker-entrypoint.sh generate -c /local/scripts/generator-cfg.yaml

# not used normally but redefine entrypoint for the case of checking this stage results
ENTRYPOINT ["bash"]

############################################
FROM golang:latest AS build

WORKDIR /app

# download dependencies - low frequency of changes
COPY go.mod go.sum ./
RUN go mod download

# copy sources - higher frequency of changes 
COPY internal/ internal/
COPY cmd/ cmd/
COPY --from=api /local/ ./

# ensure tests are passing
RUN go test ./...

# create executable - incident-log-webapi-srv
RUN CGO_ENABLED=0 GOOS=linux \
    go build \
    -ldflags="-w -s" \
    -installsuffix 'static' \
    -o ./incident-log-webapi-srv ./cmd/incident-log-api-service

############################################
FROM scratch

LABEL org.opencontainers.image.authors="Serhii Kuklych"
LABEL org.opencontainers.image.title="IncidentLog WebAPI Service"
LABEL org.opencontainers.image.description="WebAPI for hospital incident and security event management"

ENV INCIDENT_LOG_API_ENVIRONMENT=production
ENV INCIDENT_LOG_API_PORT=8080
ENV INCIDENT_LOG_API_MONGODB_HOST=mongo
ENV INCIDENT_LOG_API_MONGODB_PORT=27017
ENV INCIDENT_LOG_API_MONGODB_DATABASE=incident-log
ENV INCIDENT_LOG_API_MONGODB_COLLECTION=incident-log
ENV INCIDENT_LOG_API_MONGODB_USERNAME=root
ENV INCIDENT_LOG_API_MONGODB_PASSWORD=
ENV INCIDENT_LOG_API_MONGODB_TIMEOUT_SECONDS=5

COPY --from=build /app/incident-log-webapi-srv ./

EXPOSE 8080

ENTRYPOINT ["./incident-log-webapi-srv"]