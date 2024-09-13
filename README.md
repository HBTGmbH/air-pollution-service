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
| List all emissions of all countries for each year | GET         | /emissions/year           |
| List all emissions of all years for each country  | GET         | /emissions/country        |
| List all emissions of a specific country          | GET         | /emissions/country/{name} |

## 🚀 How to run

### Locally
 * Copy the raw data to a file, e.g. `/data/air-pollution.csv`.
 * Point env variable `AIR_POLLUTION_FILE` to the file path. 
 * Run the service:
```bash
go build cmd/server/main.go
go go run cmd/server/main.go
```