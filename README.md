# Real-Time Fraud Detection in Cryptocurrency Trading

This project implements a real-time fraud detection system for cryptocurrency trading using Go. The system fetches price data from the different sources, for instance Binance, CoinGecko, builds an isolation forest for anomaly detection, and reports detected anomalies along with statistics. The system is designed to run continuously, fetching new data at regular intervals and processing it to detect potential fraudulent activities.

## Features

- **Enabled Sources Configuration**: Perform operations based on enabled sources.
- **Real-time Data Fetching**: Continuously fetches enabled source cryptocurrency price data, e.g., Binance.
- **Isolation Forest**: Implements an isolation forest for anomaly detection.
- **Anomaly Detection**: Detects anomalies in the price data and reports them.
- **Statistics Reporting**: Provides statistics on the total number of items and anomalies detected.

## Installation

To install and run this project, follow these steps:

1. **Clone the repository:**
   ```sh
   git clone https://github.com/mdshahjahanmiah/trading-fraud-detection.git
   cd trading-fraud-detection
      ```

2. **Run and Test the application:**:
   ```sh
   go mod tidy
   go run cmd/main.go
   ```

## Usage
The main entry point of the application is the main.go file. The application fetches price data, builds an isolation forest, detects anomalies, and reports them through channels.

## Project Structure
- `main.go`: Entry point of the application.
- `source/source.go`: Contains the source e.g., BinanceSource struct and methods for fetching price data.
- `procedure/isolation_forest.go`: Contains functions for building the isolation forest.
- `detect/anomalies.go`: Contains functions for detecting anomalies.

## Implementation Details
### Isolation Forest
An isolation forest is an ensemble-based anomaly detection method. It isolates observations by randomly selecting a feature and then randomly selecting a split value between the maximum and minimum values of the selected feature.

### Anomaly Detection
Anomalies are detected based on the anomaly score. The score is calculated using the path length from the root node to the terminating node. A higher score indicates a higher likelihood of the point being an anomaly.