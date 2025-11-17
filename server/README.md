# Pasword Mango

Pasword Mango is a password manager backend designed to help you manage your credentials effortlessly. It provides a secure, high-performance RESTful API built with Go, using Google Cloud Firestore for persistent, encrypted storage. This backend is intended to be used with a frontend client, such as the planned C++ Qt desktop application.

## ✨ Features

- **Secure Credential Storage**: Passwords are encrypted at rest using **AES-256-GCM** before being stored in the database, ensuring they are never saved in plaintext.
- **RESTful API**: A clean and simple API for all CRUD (Create, Read, Update, Delete) operations on credentials.
- **High-Performance Backend**: The backend is written in [Go](https://go.dev/), ensuring speed, reliability, and excellent concurrency support.
- **Flexible Site Lookup**: API endpoints for retrieving, updating, and deleting credentials can find sites with or without a `.com` suffix, improving user experience.
- **Scalable Database**: Leverages **Google Cloud Firestore** for a scalable, serverless NoSQL database solution.

## 🛠️ Tech Stack

- **Backend**: Go
- **Frontend**: C++ with Qt
- **Database**: Google Cloud Firestore

### Justification of Technology Choices

- **Go (Backend)**: Chosen for its high performance, strong concurrency model, and straightforward deployment. Its robust standard library is ideal for building efficient and reliable network services.
- **Google Cloud Firestore (Database)**: Replaced Convex to utilize a mature, highly scalable, serverless NoSQL database from a major cloud provider. Firestore's official Go SDK provides seamless integration, powerful querying capabilities, and a generous free tier suitable for development and small-scale applications.
- **No Authentication Layer**: Authentication (like Clerk) has been removed from the backend's scope. This simplifies the architecture, positioning the service as a backend for a local-first desktop application where authentication could be handled by the OS or a local master password, rather than a multi-user web service.

## 📋 API Specification

The API provides the following endpoints for managing credentials:

### `POST /credentials`

- **Action**: Creates a new credential.
- **Body**: `{"site": "example.com", "username": "user", "password": "pw"}`
- **Responses**:
  - `201 Created`: On success.
  - `409 Conflict`: If a credential for the site already exists.
  - `400 Bad Request`: For invalid or incomplete request body.

### `GET /credentials`

- **Action**: Retrieves a list of all stored credentials.
- **Response**: `200 OK` with a JSON array of credentials.

### `GET /credentials/{site}`

- **Action**: Retrieves the credentials for a specific site. The lookup is flexible and will find `example` or `example.com`.
- **Response**: `200 OK` with the credential's JSON object.

### `PUT /credentials/{site}`

- **Action**: Updates the credentials for a specific site.
- **Body**: `{"username": "new_user", "password": "new_pw"}`
- **Response**: `200 OK` on success.

### `DELETE /credentials/{site}`

- **Action**: Deletes the credentials for a specific site.
- **Response**: `200 OK` on success.

## ⚠️ Limitations

- **No User Authentication**: The API is currently open and does not differentiate between users. It is designed for a single-user context, such as a local desktop application.
- **Single Static Encryption Key**: The application uses a single AES key loaded from an environment variable. In a production multi-user system, per-user keys derived from a master password would be more secure.
- **HTTP Only**: The server runs on HTTP and is intended for local development. For production use, it should be run behind a reverse proxy that provides TLS/SSL encryption (HTTPS).
- **No Pagination**: The `GET /credentials` endpoint retrieves all stored credentials in a single request. This will not scale well and could lead to performance issues if the database contains a large number of entries.
- **Basic Input Validation**: While the API checks for empty fields and length limits, it does not perform stricter validation on the format of inputs (e.g., ensuring a site is a valid domain name).
- **Limited Error Granularity**: For some operations like `GET /credentials/{site}`, a failure to find the site is indistinguishable from a failure to decrypt its password. Both scenarios result in a `404 Not Found` response, which can hide underlying data corruption issues.

## 🚀 Getting Started

Follow these instructions to get a local copy up and running for development and testing purposes.

### Prerequisites

- Go (version 1.21 or newer)
- Qt and a C++ compiler (e.g., GCC, Clang, MSVC)
- A Google Cloud Platform (GCP) Project with Firestore enabled.
- A GCP Service Account key (`adminkey.json`).

### Installation & Setup

1.  **Clone the repository:**

    ```sh
    git clone https://github.com/rihts-4/pasword-mango.git
    cd pasword-mango
    ```

2.  **Backend (Go):**

    - Create a `.env` file in the root directory and configure it with your GCP project details and a secure encryption key. You can generate a key with `openssl rand -hex 32`.
      ```env
      # .env
      PROJECT_ID="your-gcp-project-id"
      GOOGLE_APPLICATION_CREDENTIALS="adminkey.json"
      ENCRYPTION_KEY="your_64_character_hex_encryption_key"
      ```
    - Place your downloaded GCP service account key in the root directory and name it `adminkey.json`.
    - Install dependencies and run the server:
      ```sh
      go mod tidy
      go run .
      ```
    - The server will start on `http://localhost:8080`.

3.  **Frontend (C++ & Qt):**
    - Navigate to the `frontend` directory.
    - Use Qt Creator or your preferred C++ IDE to open the `.pro` or `CMakeLists.txt` file.
    - Build and run the project. The application will connect to the local Go backend.
