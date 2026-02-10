[![Go Reference](https://pkg.go.dev/badge/github.com/HBTGmbH/air-pollution-service.svg)](https://pkg.go.dev/github.com/HBTGmbH/air-pollution-service) 
[![Go Report Card](https://goreportcard.com/badge/github.com/HBTGmbH/air-pollution-service)](https://goreportcard.com/report/github.com/HBTGmbH/air-pollution-service) 
[![build workflow](https://github.com/HBTGmbH/air-pollution-service/actions/workflows/build.yml/badge.svg)](https://github.com/HBTGmbH/air-pollution-service/actions/workflows/build.yml)

# Air Pollution Service 🌱
A simple Go service that provides a RESTful API for retrieving emission and air pollution data for 
all countries over the last 300 years. This microservice serves as a backend service that handles requests to 
fetch historical environmental data, including CO₂ emissions, air pollutants, and other 
environmental factors, for every country around the globe. 

The required raw data can be downloaded here https://www.kaggle.com/datasets/rejeph/air-pollution.

## 💡 Endpoints

| Name                                              | HTTP Method | Route                   |
|---------------------------------------------------|-------------|-------------------------|
| OpenAPI Documentation                             | GET         | /swagger/index.html     |
| List all countries                                | GET         | /countries              |
| Returns a single country                          | GET         | /countries/{id}         |
| List all emissions of all countries for each year | GET         | /emissions/year/        |
| Get all emissions of a specific year              | GET         | /emissions/year/{id}    |
| List all emissions of all years for each country  | GET         | /emissions/country/     |
| Get all emissions of a specific country           | GET         | /emissions/country/{id} |
| Prometheus Metrics                                | GET         | /metrics                |

## 🚀 How to run

### Environment

| Name                            | Default           | Route                                                                                                                                                                                                                                                                                                                                                                       |
|---------------------------------|-------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| SERVER_PORT                     | 8080              | Port used by the service.                                                                                                                                                                                                                                                                                                                                                   |
| AIR_POLLUTION_FILE              | air-pollution.csv | Path to CSV file with [raw data](ttps://www.kaggle.com/datasets/rejeph/air-pollution).                                                                                                                                                                                                                                                                                      |
| SERVER_IDLE_TIMEOUT             | 60s               | Idle timeout for HTTP server. The server will keep idle client connections open for this duration.                                                                                                                                                                                                                                                                          |
| SERVER_SHUTDOWN_SLEEP_DURATION  | 0s                | Duration to sleep before initiating server shutdown. This should be set to a duration after which no new connections are expected. This should always be used when deploying in a Kubernetes environment, to allow kube-proxy to update iptables/netfilter rules, which usually take ~5s. A value of `15s` is advised.                                                      |
| SERVER_DRAIN_TIMEOUT            | 30s               | Duration to wait for existing (idle) connections to be closed by the clients. During this time, the server will still accept new connections and requests but respond with 'Connection: close' header to all requests. This should always be set to the maximum client-side idle timeout, to avoid the server closing a connection right as the client sends a new request. |
| SERVER_SHUTDOWN_TIMEOUT         | 30s               | Duration to wait for the server to shutdown gracefully. During this time, the server will not accept any more connections and requests, but wait for existing requests to finish. When not using `SERVER_DRAIN_TIMEOUT` this should be set to the maximum request duration.                                                                                                 |

### Locally
 * Download the raw data to a file, e.g. `/data/air-pollution.csv`.
 * Point env variable `AIR_POLLUTION_FILE` to the file path. 
 * Run the service:

#### Without cloning the Git repository
When the Go SDK is installed, the service can be installed and run with the following commands:
```bash
go install github.com/HBTGmbH/air-pollution-service@latest
air-pollution-service
```
The binary `air-pollution-service` is now available in the `$GOPATH/bin` directory, such as `~/go/bin/` by default.

#### With cloning the Git repository
When working with the source code, the service can be run with the following command:
```bash
go run main.go
```

### Docker
* Download the raw data to a file, e.g. `/data/air-pollution.csv`.
```bash
docker pull ghcr.io/hbtgmbh/air-pollution-service:latest
docker run -v /data:/data -p 8080:8080 -e AIR_POLLUTION_FILE=/data/air-pollution.csv ghcr.io/hbtgmbh/air-pollution-service:latest
```

## 🔨 How to build

### Using Go
```bash
go build -o server
```

### Using Docker
```bash
docker build -t air-pollution-service .
```

### Using Make
```bash
make build
```

### Swagger API Documentation
For the Swagger API documentation, the project https://github.com/swaggo/swag is used. This generates the Swagger API documentation and integrates with a web framework (as such "Chi" used here) to render the Swagger Web UI at runtime.

The sources of this api-pollution-service come with an up-to-date generated Swagger documentation to facilitate easier Go compilation after checking out this project, as well as to make use of Git diffs between commits to review changes in the generated files.

To (re-)generate the Swagger API documentation, the `swag` CLI command must be installed:
```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

After that, the documentation inside of the `docs/` directory can be (re-)generated by running:
```bash
swag init
```

## 📝 TODO
Authentication: https://github.com/HBTGmbH/air-pollution-service/issues/3
