[![Go Reference](https://pkg.go.dev/badge/github.com/HBTGmbH/air-pollution-service.svg)](https://pkg.go.dev/github.com/HBTGmbH/air-pollution-service) 
[![Go Report Card](https://goreportcard.com/badge/github.com/HBTGmbH/air-pollution-service)](https://goreportcard.com/report/github.com/HBTGmbH/air-pollution-service) 
[![build workflow](https://github.com/HBTGmbH/air-pollution-service/actions/workflows/docker-build.yml/badge.svg)](https://github.com/HBTGmbH/air-pollution-service/actions/workflows/docker-build.yml)

# Air Pollution Service 🌱
A simple Go service that provides a RESTful API for retrieving emission and air pollution data for 
all countries over the last 300 years. This microservice serves as a backend service that handles requests to 
fetch historical environmental data, including CO₂ emissions, air pollutants, and other 
environmental factors, for every country around the globe. 

The required raw data can be downloaded here https://www.kaggle.com/datasets/rejeph/air-pollution.

## 💡 Endpoints

| Name                                              | HTTP Method | Route                     |
|---------------------------------------------------|-------------|---------------------------|
| List all countries                                | GET         | /countries                |
| Returns a single country                          | GET         | /countries/{name}         |
| List all emissions of all countries for each year | GET         | /emissions/year/          |
| Get all emissions of a specific year              | GET         | /emissions/year/{year}    |
| List all emissions of all years for each country  | GET         | /emissions/country/       |
| Get all emissions of a specific country           | GET         | /emissions/country/{name} |

## 🚀 How to run

### Environment

| Name               | Default           | Route                                                                                  |
|--------------------|-------------------|----------------------------------------------------------------------------------------|
| SERVER_PORT        | 8080              | Port used by the service.                                                              |
| AIR_POLLUTION_FILE | air-pollution.csv | Path to CSV file with [raw data](ttps://www.kaggle.com/datasets/rejeph/air-pollution). |

### Locally
 * Download the raw data to a file, e.g. `/data/air-pollution.csv`.
 * Point env variable `AIR_POLLUTION_FILE` to the file path. 
 * Run the service:
```bash
go build cmd/server/main.go
go go run cmd/server/main.go
```

### Docker
* Download the raw data to a file, e.g. `/data/air-pollution.csv`.
```bash
docker pull ghcr.io/hbtgmbh/air-pollution-service:20
docker run --mount type=bind,src=/dat,dst=/data --publish 8080:8080 --env AIR_POLLUTION_FILE=/data/air-pollution.csv ghcr.io/hbtgmbh/air-pollution-service:20
```

## 📝 TODO
Authentication: https://github.com/HBTGmbH/air-pollution-service/issues/3