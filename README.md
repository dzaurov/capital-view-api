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
    go install [github.com/swaggo/swag/cmd/swag@latest](https://www.google.com/search?q=https://github.com/swaggo/swag/cmd/swag%40latest)
    # Ensure your $GOPATH/bin or $HOME/go/bin is in your system's PATH
    ```

## Setup and Installation

Follow these steps to get the project running locally:

1.  **Clone the Repository:**
    Replace `<your-github-username>` and `<your-repo-name>` with the actual GitHub user/organization and repository name.
    ```bash
    git clone [https://github.com/](https://github.com/)<your-github-username
