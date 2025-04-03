# Go Company Data API

This project provides a RESTful API built with Go (using Gin and GORM) to manage and search company registration data, member information, beneficial owners, and financial statements stored in a SQLite database. API documentation is generated using Swagger.

## Prerequisites

Before you begin, ensure you have the following installed:

1.  **Go:** Version 1.18 or higher. ([Download Go](https://golang.org/dl/))
2.  **Git:** Required to clone the repository. ([Download Git](https://git-scm.com/downloads))
3.  **SQLite3 Development Libraries:** The Go SQLite driver (`mattn/go-sqlite3`) requires CGO, which needs the SQLite3 development libraries.
    * **Debian/Ubuntu:** `sudo apt-get update && sudo apt-get install libsqlite3-dev`
    * **Fedora/CentOS/RHEL:** `sudo dnf install libsqlite3-devel` or `sudo yum install libsqlite3-devel`
    * **macOS (using Homebrew):** `brew install sqlite` (usually includes necessary headers)
    * **Windows:** Setting up CGO on Windows can be complex. Consider using a pre-compiled SQLite or setting up a GCC environment like MinGW-w64 or TDM-GCC. Refer to `mattn/go-sqlite3` documentation for details.
4.  **Swag CLI:** Used to generate Swagger documentation from code annotations.
    ```bash
    go install [github.com/swaggo/swag/cmd/swag@latest]
    # Ensure your $GOPATH/bin or $HOME/go/bin is in your system's PATH
    ```

## Setup and Installation

Follow these steps to get the project running locally:

1.  **Clone the Repository:**
    Replace `<your-github-username>` and `<your-repo-name>` with the actual GitHub user/organization and repository name.
    ```bash
    git clone [https://github.com/](https://github.com/)<your-github-username>/<your-repo-name>.git
    ```

2.  **Navigate to Project Directory:**
    ```bash
    cd <your-repo-name>
    ```

3.  **Install Go Dependencies:**
    This command reads the `go.mod` file and installs the required Go modules listed in `go.sum`.
    ```bash
    go mod tidy
    # or alternatively: go mod download
    ```

4.  **Generate Swagger Documentation Files:**
    This command parses the Swagger annotations in your Go code (`//@...`) and generates the necessary files in the `docs/` directory.
    ```bash
    swag init
    ```

## Running the Application

1.  **Start the API Server:**
    Execute the following command from the project's root directory:
    ```bash
    go run main.go
    ```
    You should see output indicating the server is starting, typically listening on port 8080:
    ```
    [GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

    [GIN-debug] Listening and serving HTTP on :8080
    INFO[0000] Database connection established
    INFO[0000] Database migrated
    INFO[0000] Server starting on port :8080
    ```
    The server needs to remain running in this terminal window to handle API requests.

## Accessing the API Documentation (Swagger UI)

Once the application is running:

1.  Open your web browser.
2.  Navigate to: `http://localhost:8080/swagger/index.html`

This page provides interactive documentation for all available API endpoints. You can view details, models, and execute test requests directly from the browser.

## Database

* The application uses **SQLite** as its database.
* The database file is named `mydata.db` and will be created in the project's root directory.
* The application uses GORM's `AutoMigrate` feature. On the first run (or if the `mydata.db` file is missing), it will automatically create the database file and all the necessary tables based on the models defined in the `models/` directory.

---

You can save this content as `README.md` in the root of your GitHub repository. Remember to replace the placeholder clone URL with the correct one.
