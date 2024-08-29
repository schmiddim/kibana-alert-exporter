# Kibana Alert Exporter

[![codecov](https://codecov.io/gh/schmiddim/kibana-alert-exporter/graph/badge.svg?token=yzZKPtqT4e)](https://codecov.io/gh/schmiddim/kibana-alert-exporter)
[![Docker Pulls](https://img.shields.io/docker/pulls/schmiddim/kibana-alert-exporter.svg)](https://hub.docker.com/r/schmiddim/kibana-alert-exporter/)

# Overview

Kibana Alert Exporter is a tool designed to export alerts from Kibana. It provides a simple interface to interact with
Kibana's alerting API and retrieve alert data.

# Features

- Export alerts from Kibana - Supports multiple alert types - Easy integration with CI/CD pipelines

# Installation

## Using Docker

You can pull the Docker image from Docker Hub: 
```sh 
docker pull schmiddim/kibana-alert-exporter 
```

## From Source
To build and run the project from source, ensure you have Go installed and run the following commands:

```sh 
git clone https://github.com/schmiddim/kibana-alert-exporter.git 
cd kibana-alert-exporter 
go build -o kibana-alert-exporter 
./kibana-alert-exporter 

```  
# Usage
## Command Line
Run the exporter with the following command:
```sh 
export KIBANA_URL=<KIBANA_URL>
exporter API_KEY=<API_KEY>
./kibana-alert-exporter --kibana-url <KIBANA_URL> --api-key <API_KEY> 
```  
## Docker
Run the Docker container with the necessary environment variables:
```sh 
docker run -e KIBANA_URL=<KIBANA_URL> -e API_KEY=<API_KEY> schmiddim/kibana-alert-exporter 
```  
## Helm Chart
To deploy the exporter using Helm, use the following commands:
```sh
helm repo add schmiddim https://schmiddim.github.io/helm-charts/
helm repo update
helm upgrade kibana-alert-exporter schmiddim/kibana-alert-exporter --install
```
# Configuration
The following environment variables can be set to configure the exporter:  
- `KIBANA_URL`: The URL of the Kibana instance.
- `API_KEY`: The API key for authenticating with Kibana.  
# Development
## Running Tests
To run the tests, use the following command: 
```sh 
go test ./... 
```  
## Contributing
Contributions are welcome! Please open an issue or submit a pull request for any changes.  
## License
This project is licensed under the MIT License. See the `LICENSE` file for details.