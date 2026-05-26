param (
    $command
)

if (-not $command) {
    $command = "start"
}

$ProjectRoot = "${PSScriptRoot}/.."

$env:INCIDENT_LOG_API_ENVIRONMENT = "Development"
$env:INCIDENT_LOG_API_PORT = "8080"

switch ($command) {
    "start" {
        go run ${ProjectRoot}/cmd/incident-log-api-service
    }
    "openapi" {
        docker run --rm -ti -v ${ProjectRoot}:/local openapitools/openapi-generator-cli generate -c /local/scripts/generator-cfg.yaml
    }
    default {
        throw "Unknown command: $command"
    }
}