# automation-backend

> High-performance Go backend for orchestrating and managing Python-based automation and simulation tasks.

##  Status Badges

| Build/Test | Coverage | Go Version |
| :---: | :---: | :---: |
| [![Build Status](https://img.shields.io/badge/build-passing-brightgreen)](https://github.com/CBYeuler/automation-backend/actions) | [![Coverage](https://img.shields.io/badge/coverage-85%25-yellowgreen)](YOUR_COVERAGE_REPORT_LINK) | [![Go Version](https://img.shields.io/badge/Go-1.21+-blue)](https://go.dev/) |

##  GitHub Topics/Tags

`Go`, `Gin`, `Gorm`, `Python`, `Simulation`, `REST-API`, `Automation`

##  Project Overview

### What is this project?

This project, `automation-backend`, is a high-performance orchestration system designed to manage and execute complex automation and simulation workflows. It functions as a robust **REST-API** gateway, enabling external systems (like web frontends or scheduling services) to trigger, monitor, and retrieve results from computationally intensive tasks.

### Where can it be used?

It is ideally suited for:
* **Continuous Integration/Deployment (CI/CD):** Running automated performance, load, or functional tests as part of a pipeline.
* **Financial Modeling/Scientific Computing:** Managing batches of complex simulations where coordination and data logging are critical.
* **Digital Twin Systems:** Orchestrating simulations that model real-world processes or physical infrastructure.

### What problem does it solve?

The primary problem it solves is the need for a reliable, scalable, and concurrent platform to run long-running, resource-heavy automation or simulation tasks. This backend provides essential services like state persistence, request queuing, concurrent task handling, and standardized result reporting, ensuring stability and continuous operation without manual oversight.

## Tech & Design Decisions

### Why Go for concurrency?

Go was selected specifically for its superior **concurrency model** using goroutines. This is vital for a backend that must handle many simultaneous incoming requests, efficiently manage long-running background tasks, and maintain high throughput without heavy system resource consumption. Gin provides a fast API framework, and Gorm handles reliable, structured database interaction and state persistence.

### Why Python for the simulator?

Python is used for the actual simulation logic because it boasts a mature and extensive ecosystem of scientific, data analysis, and specialized simulation libraries. It is the ideal language for rapid development of complex algorithms, while Go remains the high-performance *orchestrator* that calls the Python components.

##  Installation

### Prerequisites

* Go (version 1.21 or later)
* Python (version 3.8 or later)
* A running database instance (PostgreSQL/SQLite).

### Getting Started

1.  **Clone the repository:**
    ```bash
    git clone [https://github.com/CBYeuler/automation-backend.git](https://github.com/CBYeuler/automation-backend.git)
    cd automation-backend
    ```
2.  **Install Go dependencies:**
    ```bash
    go mod download
    ```
3.  **Install Python dependencies (for simulator):**
    ```bash
    pip install -r simulator/requirements.txt
    ```

##  Usage with Makefile

The project uses a `Makefile` to simplify common development tasks:

| Command | Description |
| :---: | :---: |
| `make run` | Builds the Go binary and starts the server. |
| `make test` | Runs all Go unit and integration tests. |
| `make seed` | Executes the database seed script to populate initial data. |

To run the application:
```bash
make run
```
### TODO List

- Implement database migration system (e.g., using golang-migrate).

- Implement a proper job queuing mechanism (e.g., integrating with Redis or Kafka).

- Dockerize the application for easier deployment and portability.

- Add OpenAPI/Swagger documentation for the REST API endpoints.


```text
MIT License

Copyright (c) 2025 CBYeuler
```
