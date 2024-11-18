# WasteTags

WasteTags is a web application for generating chemical waste labels with embedded QR codes. It allows users to input chemical information, generates a unique waste tag, and produces a label that can be printed and attached to waste containers. The application uses a SQLite database to retrieve chemical component data and generates QR codes containing waste information for easy scanning and tracking.

## Future Development Ideas

- Investigate using [WAILS](https://wails.io/) for the front end. Can be dockerized & build Windows, macOS, and Linux desktop apps & allows for Go backend
-

## Table of Contents

- [Features](#features)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Usage](#usage)
- [Makefile Commands](#makefile-commands)
- [Docker Deployment](#docker-deployment)
- [Project Structure](#project-structure)

## Features

- Generate unique waste tags using UUIDs.
- Retrieve chemical component data from a SQLite database.
- Generate QR codes containing waste information.
- Render labels using HTML templates.
- Embed static assets using Go's `embed` package.
- Run the application locally or inside a Docker container.

## Prerequisites

- [Go (Golang) 1.17 or later](https://golang.org/dl/)
- [Git](https://git-scm.com/downloads)
- [Docker](https://www.docker.com/products/docker-desktop) (for Docker deployment)
- [Docker Compose](https://docs.docker.com/compose/install/) (for Docker deployment)

## Installation

Clone the repository and navigate to the project directory:

```sh
git clone https://github.com/yourusername/wastetags.git
cd wastetags
```

## Usage

### Running Locally

Ensure that you have Go installed and your `GOPATH` is set up.

1. **Install Dependencies**

   ```sh
   make deps
   ```

2. **Build the Application**

   ```sh
   make build
   ```

3. **Run the Application**

   ```sh
   make run
   ```

   The application will start on `http://localhost:8080`.

4. **Access the Application**

   Open your web browser and navigate to `http://localhost:8080` to use the waste label generator.

### Running with Docker

Ensure Docker and Docker Compose are installed.

1. **Build and Start the Docker Container**

   ```sh
   make docker-up
   ```

2. **Access the Application**

   Open your web browser and navigate to `http://localhost:8080`.

3. **Stop the Docker Container**

   ```sh
   make docker-down
   ```

## Makefile Commands

The project includes a `Makefile` to streamline common tasks. Below is an explanation of the available commands:

- **Build Commands**

  - `make` or `make build`: Builds the Go application and outputs the binary to the `bin/` directory.

    ```sh
    make build
    ```

- **Run Commands**

  - `make run`: Builds the application (if not already built) and runs it.

    ```sh
    make run
    ```

- **Clean Commands**

  - `make clean`: Removes the built binaries and cleans up the `bin/` directory.

    ```sh
    make clean
    ```

- **Dependency Management**

  - `make deps`: Ensures all module dependencies are up to date using `go mod tidy`.

    ```sh
    make deps
    ```

- **Testing Commands**

  - `make test`: Runs all tests in the project.

    ```sh
    make test
    ```

- **Code Quality Commands**

  - `make fmt`: Formats the code using `go fmt`.

    ```sh
    make fmt
    ```

  - `make lint`: Lints the code using `golangci-lint` (requires installation).

    ```sh
    make lint
    ```

- **Code Generation**

  - `make generate`: Runs code generation tools (if applicable).

    ```sh
    make generate
    ```

- **Docker Commands**

  - `make docker-up`: Builds the Docker image and starts the services defined in `docker-compose.yml`.

    ```sh
    make docker-up
    ```

  - `make docker-down`: Stops and removes the Docker services.

    ```sh
    make docker-down
    ```

  - `make docker-build-no-cache`: Builds the Docker image without using the cache.

    ```sh
    make docker-build-no-cache
    ```

  - `make docker-logs`: Follows the logs from the Docker services.

    ```sh
    make docker-logs
    ```

- **Help Command**

  - `make help`: Displays the help message with a list of available commands.

    ```sh
    make help
    ```

## Docker Deployment

The application can be containerized using Docker for ease of deployment.

### Build and Run with Docker Compose

1. **Start the Services**

   ```sh
   make docker-up
   ```

2. **Access the Application**

   Visit `http://localhost:8080` in your web browser.

3. **View Logs**

   ```sh
   make docker-logs
   ```

4. **Stop the Services**

   ```sh
   make docker-down
   ```

### Notes

- The `docker-compose.yml` file is located in the `deployments/` directory.
- The `data/` directory is mounted as a volume inside the Docker container to persist the SQLite database.

## Project Structure

```txt
wastetags
├── Makefile
├── bin
│   └── wastetags
├── build
│   └── package
│       └── Dockerfile
├── cmd
│   └── wastetags
│       ├── main.go
│       └── templates
│           ├── tag-form.html
│           └── tag.html
├── data
│   └── chemicals.sqlite3
├── deployments
│   └── docker-compose.yml
├── go.mod
├── go.sum
└── pkg
    ├── database
    │   └── database.go
    └── qrcodegen
        └── qrcodegen.go
```

- `cmd/wastetags/main.go`: The entry point of the application.
- `cmd/wastetags/templates/`: HTML templates for rendering forms and labels.
- `pkg/database/`: Package handling database operations.
- `pkg/qrcodegen/`: Package for generating QR codes.
- `data/chemicals.sqlite3`: SQLite database containing chemical information.
- `build/package/Dockerfile`: Dockerfile for building the Docker image.
- `deployments/docker-compose.yml`: Docker Compose configuration.
- `Makefile`: Build and management commands.
