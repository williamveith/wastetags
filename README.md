# Waste Tagging System

The Waste Tagging System is a web application designed to streamline the process of creating and managing waste labels. It provides a user-friendly interface for adding new chemicals, creating mixtures, and generating QR-code-enhanced labels for containers containing hazardous or non-hazardous waste. This system helps ensure regulatory compliance, improves inventory tracking, and simplifies waste disposal management workflows.

## Features

- **Add New Chemical**: Enter CAS numbers and chemical names through a guided form. The application validates CAS numbers to ensure data integrity.
- **Add Mixture**: Combine individual chemicals to form a mixture. Each mixture entry can be reused when generating waste labels.
- **Generate Waste Labels**: Create and print waste tags that include:
  - Unique QR codes for each label.
  - Information such as room number, generator name, chemical components, and their percentages.
  - Multiple copies of labels with unique identifiers generated on the fly.
- **Embedded Resources**: Comes with embedded HTML templates, CSS, and SQL schema. This makes deployment and configuration simpler.
- **Scalable API**: Integrated endpoints for managing database entries and generating dynamic QR codes for labeling.

## Technology Stack

- **Backend**: Go (Golang) with [Gin](https://github.com/gin-gonic/gin) for routing and HTTP request handling.
- **Database**: SQLite3 for local storage of chemicals, mixtures, locations, containers, units, and states.
- **QR Code Generation**: Leveraging `github.com/williamveith/wastetags/pkg/qrcodegen` for creating data URIs from JSON payloads.
- **Templates & Assets**: Embedded using Go’s `embed` package, making the application self-contained.

## Prerequisites

- **Go (1.19+ recommended)**: Ensure [Go](https://go.dev/) is installed.
- **SQLite3**: The database is automatically initialized from embedded SQL files. No separate installation required unless you want to inspect the database independently.
- **Git**: If you plan on cloning the repository directly.

## Build The App

1. Clone the Repository
2. Navigate to Project Directory
3. Build the Project

  ```sh
  git clone https://github.com/yourusername/wastetags.git
  cd wastetags
  make build
  ```

## Use The App

Double click the binary generated in the build located in the bin directory or use the make run command in the wastetags directory. The UI can be accessed at `http://localhost:8080`

From here, you can:

- Navigate to **Add Chemical**: `/addchemical`
- Navigate to **Add Mixture**: `/add-mixture`
- Navigate to **Create Tag**: `/create-tag`

## Feature Details

1. **Adding a Chemical**:
   - Go to `/addchemical`.
   - Enter the CAS number in three parts (e.g., `64-17-5` for ethanol).
   - Provide the chemical name.
   - Click **Add Chemical**. The CAS number will be validated before submission.

2. **Adding a Mixture**:
   - Go to `/add-mixture`.
   - Specify a CAS number and name for your mixture (e.g., a solvent mixture).
   - After submission, the mixture information can be used when creating waste tags.

3. **Creating a Waste Tag**:
   - Go to `/create-tag`.
   - Select the **Location**, **Chemical Name**, container details, and amount.
   - Submit the form to generate a waste label with a QR code.
   - Print the label by using your browser’s print functionality. You can also specify multiple copies and generate unique QR codes for each one.

## Project Structure

- **`cmd/wastetags`**: Main application entry point.
- **`pkg/database`**: Database interaction and initialization.
- **`templates/`**: Embedded HTML templates for the user interface.
- **`assets/`**: Embedded CSS styles and images for UI styling.
- **`query/`**: Embedded SQL files for schema setup and queries.
- **`Makefile`**: Build and run targets for the project.
- **`Dockerfile` and `docker-compose.yml` (if present)**: For containerized deployment (adjust as needed).

## Customization & Configuration

- **Port**: The application runs by default on port `:8080`. Modify the `r.Run(":8080")` line in `runLabelMaker()` to change the port.
- **Database Path**: The SQLite database defaults to `data/chemicals.sqlite3`. Modify the `init()` function in the code to point to a different database file if desired.

## Troubleshooting

- **Schema Load Errors**: Ensure that the `query/schema.sql` file is present and can be read. The database initialization requires it.
- **Permission Errors**: If running on a restricted environment, ensure the application has the necessary permissions to read embedded files and create the database file.
- **CORS Issues**: Since the application is typically accessed from the same origin, CORS settings aren’t commonly needed. If deploying behind a reverse proxy or integrating with another frontend, you may need to configure CORS policies within Gin.
