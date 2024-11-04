# Receipt Processor

## Objective

The objective is to create a web service that processes receipt data and calculates points based on specific rules. The service should expose two endpoints: one for processing receipts and another for retrieving points for a processed receipt.

## Getting Started

### Prerequisites

The application needs following installed.

- [Go](https://go.dev/doc/install) (v1.23.2+)
- [Docker](https://docs.docker.com/engine/install/)
- [swag](https://github.com/swaggo/swag)

### Usage: To Run the application

#### 1. User docker
Make sure that the docker is up and running. Then

- Build
    ```bash
    docker compose build
    ```
- Run
    ```bash
    docker compose up
    ```

#### 2. User Go

- Install dependency
    ``` bash
    go mod download
    ```
- Run Test
    ```bash
    go test ./...
    ```
- Run Swag
    ```bash
    swag init -g cmd/server/main.go --parseDependency --parseInternal
    ```
- Run Server
    ```bash
    air
    ```
    or
    ```bash
    go run cmd/server/main.go
    ```

#### 3. Makefile
- Run test, init swagger and then build
    ```bash
    make all
    ```
- Run Server
    ```bash
    make run
    ```

### Usage: To Test the endpoints

The application will run on `http://localhost:9080` port

1. Using Swagger api docs runs on [http://localhost:9080/swagger/index.html](http://localhost:9080/swagger/index.html)
    Use the swagger docs to test the endpoints.
2. Using Curl
    - To add a receipt
        ```
        curl -X POST   -H "Content-Type: application/json"   -d '{
        "retailer": "Target",
        "purchaseDate": "2022-01-01",
        "purchaseTime": "13:01",
        "items": [
            {
            "shortDescription": "Mountain Dew 12PK",
            "price": "6.49"
            },{
            "shortDescription": "Emils Cheese Pizza",
            "price": "12.25"
            "price": "1.26"
            },{
            "shortDescription": "Doritos Nacho Cheese",
            "price": "3.35"
            },{
            "shortDescription": "   Klarbrunn 12-PK 12 FL OZ  ",
            "price": "12.00"
            }
        ],
        "total": "35.35"
        }' http://localhost:9080/receipt/process
        ```
    - To get the point for the receipt id
        ```
        curl http://localhost:9080/receipt/{id}/points
        ```
