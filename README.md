# Waste Tagging System

The Waste Tagging System is a web application designed to streamline the process of creating and managing waste labels. It provides a user-friendly interface for adding new chemicals, creating mixtures, and generating QR-code-enhanced labels for containers containing hazardous or non-hazardous waste. This system helps ensure regulatory compliance, improves inventory tracking, and simplifies waste disposal management workflows.

## 🗂️ Table of Contents

1. [Features](#features)
2. [Requirements](#requirements)
3. [Installation & Build](#installation--build)
4. [Configuration](#configuration)
5. [Database Management](#database-management)
6. [Routes & Handlers](#routes--handlers)
7. [UI Overview](#ui-overview)
8. [Usage Flow](#usage-flow)

---

## Features

- **Add Chemicals**: Input CAS number and chemical name to store in the database.  
- **Add Mixtures**: Construct mixtures by combining multiple chemicals.  
- **Waste Tag Generation**: Gather relevant information (location, container size, etc.) and generate a printable label with a QR code.  
- **QR Code Data**: Encodes essential details of the waste tag, enabling scanning for future reference or inventory checks.  
- **Embedded Assets**: All CSS, JavaScript, images, and SQL scripts are embedded in the Go binary, simplifying deployment.

---

## Requirements

- Go **1.19+** (or a version compatible with Go modules).  
- (Optional) Protobuf tools if you plan to extend or regenerate the embedded data (but not necessary for normal usage).  
- SQLite library is used under the hood by the **[go-sqlite3](https://github.com/mattn/go-sqlite3)** driver.

---

## Installation & Build

1. **Clone the Repository:**

   ```bash
   git clone https://github.com/williamveith/wastetags.git
   cd wastetags
   ```

2. **Download Dependencies:**

   ```bash
   go mod tidy
   ```

3. **Build:**

   ```bash
   go build -o wastetags
   ```

4. **Run:**

   ```bash
   ./wastetags
   ```

   Or:

   ```bash
   go run main.go
   ```

By default, the service listens on port **:8080**. Access the application in your browser at [http://localhost:8080](http://localhost:8080).

---

## Configuration

Wastetags supports **two** ways to load config:

1. **`--config` flag**: Provide a path to a custom JSON config.  
2. **`--dev` flag**: Loads a `dev.json` config from the embedded `configs/` folder.  

These JSON configs typically include a `database_path`, which determines where the SQLite file resides. Example usage:

```bash
./wastetags --config /path/to/config.json
```

or

```bash
./wastetags --dev
```

If no flags are provided, the application attempts to load a config based on the OS (e.g., `configs/darwin.json`, `configs/windows.json`, etc.) or fails if not found.

---

## Database Management

1. **Initialization**  
   - On first run, Wastetags checks for the existence of your database file (`database_path` from config).  
   - If not found, it creates and initializes one using the embedded `schema.sql`.  
2. **Protobuf Import/Export**  
   - Some tables can be populated from `.bin` files (Protobuf data). This occurs automatically when the database is first created (if `NeedsInitialization` is true).  
3. **Schema**  
   - Found in `query/schema.sql`. Contains table definitions for **chemicals**, **mixtures**, **locations**, **containers**, **units**, **states**, etc.

---

## Routes & Handlers

Below is a summary of the key routes and their associated handlers:

| **Route**              | **Method** | **Handler**          | **Description**                                                  |
|------------------------|-----------:|----------------------|------------------------------------------------------------------|
| `/home`                | GET        | `HomePage`           | Displays the main/home page.                                     |
| `/waste-tag-form`      | GET        | `MakeWasteTagForm`   | Renders a form to gather details for generating a waste tag.     |
| `/waste-tag`           | POST       | `MakeWasteTag`       | Processes form submission and produces a printable label.         |
| `/add-chemical`        | GET, POST  | `AddChemical`        | Displays a chemical input form and inserts new records.          |
| `/add-mixture`         | GET, POST  | `AddMixture`         | Displays a mixture input form and inserts new mixtures.          |
| `/api/generate-qr-code`| POST       | `MakeNewQRCode`      | Returns a data URI for a new QR code with embedded JSON metadata. |

Each handler typically reads an embedded SQL file from the `query/` directory, interacts with the database using the methods in `pkg/database/`, and passes data to or from the front-end forms.

---

## UI Overview

All HTML templates are embedded in the Go binary (using `//go:embed`). These templates rely on partials (`component-navbar.html`, `component-footer.html`) for consistent site headers and footers.

Below is a brief look at some core templates and their functions:

### **add-chemical.html**

```html
<!DOCTYPE html>
<html lang="en">
  <head>
    ...
    <title>Add New Chemical</title>
  </head>
  <body>
    {{ template "component-navbar.html" . }}
    <h1>Add New Chemical</h1>
    <form
      id="addChemicalForm"
      class="parent"
      action="/add-chemical"
      method="POST"
      autocomplete="off">
      <div class="content-div">
        <label for="cas1">CAS Number:</label>
        <input id="cas1" ... />
        <span> - </span>
        <input id="cas2" ... />
        <span> - </span>
        <input id="cas3" ... />
      </div>
      <div class="content-div">
        <label for="chemical-name">Chemical Name:</label>
        <input id="chemical-name" name="chemical-name" required />
      </div>
      <button type="submit">Add Chemical</button>
    </form>
    {{ template "component-footer.html" . }}
    ...
  </body>
</html>
```

- **Purpose**: Gathers CAS number and chemical name.  
- **Front-end Validation**: Basic checks via JS (e.g., digit-only inputs).

### add-mixture.html

```html
<!DOCTYPE html>
<html lang="en">
  <head>
    ...
    <title>Add New Mixture</title>
  </head>
  <body>
    {{ template "component-navbar.html" . }}
    <h1>Add New Mixture</h1>
    <h2>Still needs to be developed</h2>
    {{ template "component-footer.html" . }}
    ...
  </body>
</html>
```

- **Purpose**: (Work in progress) Allows creation of a mixture from multiple components.

### **home.html**

```html
<!DOCTYPE html>
<html lang="en">
  <head>
    ...
    <title>Waste Label Maker</title>
  </head>
  <body>
    {{ template "component-navbar.html" . }}
    <h1>Welcome to the Waste Tagging System</h1>
    <p>Use the navigation above to create tags & add new chemicals and mixtures</p>
    {{ template "component-footer.html" . }}
  </body>
</html>
```

- **Purpose**: Landing page.

### **waste-tag-form.html**

```html
<!DOCTYPE html>
<html lang="en">
  <head>
    ...
    <title>Generate Waste Label</title>
  </head>
  <body>
    {{ template "component-navbar.html" . }}
    <h1>Generate Waste Label</h1>
    <form id="wasteTagForm" ...>
      ...
      <div class="div2">
        <label for="chemName">Chemical Name</label>
        <select id="chemName" name="chemName" required>
          {{ range .Components }}
          <option value="{{ .chem_name }}">{{ .chem_name }}</option>
          {{ end }}
        </select>
      </div>
      ...
      <button type="submit">Generate Label</button>
    </form>
    {{ template "component-footer.html" . }}
    ...
  </body>
</html>
```

- **Purpose**: The main user form to gather location, container type, physical state, and so on before generating a waste tag.

### **waste-tag.html**

```html
<!DOCTYPE html>
<html lang="en">
  <head>
    ...
    <title>Waste Label Maker</title>
  </head>
  <body>
    {{ template "component-navbar.html" . }}
    <div class="tag-utilities">
      <label for="numberOfCopies">Number Of Copies</label>
      <input id="numberOfCopies" type="text" value="1" />
      <button onclick="window.print()">Print Tag</button>
    </div>
    <div class="label">
      <div class="qr-code">
        <img id="qrCodeDataUri" src="{{.QrCodeBase64}}" alt="QR Code" />
      </div>
      ...
    </div>
    {{ template "component-footer.html" . }}
    ...
  </body>
</html>
```

- **Purpose**: Displays the final waste label with a QR code, offering an option to print multiple copies.

---

## Usage Flow

1. **Add Chemicals**: Navigate to **Add Chemical** → Provide CAS info and chemical name → Submit → Stored in the DB.  
2. **Add Mixtures**: Navigate to **Add Mixture** → Provide mixture components → Submit.  
3. **Generate Tag**: Navigate to **Create Tag** → Fill out the form → Submit → The system fetches details (like chemical composition) from the DB → Renders a label with a QR code and relevant info.  
4. **Print**: Choose the number of copies → Print directly from the browser.
