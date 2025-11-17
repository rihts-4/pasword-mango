# Pasword Mango

Pasword Mango is a full-stack password manager designed to help you manage your credentials effortlessly. It combines a secure, high-performance RESTful API built with Go and a native C++ Qt desktop application for a seamless user experience.

## ✨ Features

### Backend (Go)

- **Secure Credential Storage**: Passwords are encrypted at rest using **AES-256-GCM**.
- **RESTful API**: A clean and simple API for all CRUD operations.
- **High-Performance**: Written in Go for speed and reliability.
- **Flexible Site Lookup**: Finds sites with or without a `.com` suffix.
- **Scalable Database**: Leverages Google Cloud Firestore.

### Frontend (C++ & Qt)

- **Intuitive UI**: A clean and simple desktop interface built with Qt.
- **Full CRUD Functionality**: Add, view, update, and delete passwords directly from the app.
- **Secure Password Toggling**: Passwords are redacted by default. A confirmation dialog is required to view them in plain text, preventing accidental exposure.

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
- **Response**: `200 OK` on success, `404 Not Found` if the site does not exist.

## ⚠️ Limitations & Future Updates

### Current Limitations

- **No User Authentication**: The API is open and designed for a single-user, local-first context. Future versions could add a master password to encrypt the local database or the communication key.
- **Single Static Encryption Key**: The backend uses a single AES key from an environment variable. A more secure approach would involve per-user keys derived from a master password.
- **HTTP Only**: The server runs on HTTP. For production, it should be run behind a reverse proxy that provides TLS/SSL.
- **No Pagination**: The main list loads all credentials at once, which could be slow with many entries.
- **Basic UI Features**: The frontend is functional but lacks advanced features like search, sorting, or password generation.

### Future Updates

- **Copy to Clipboard**: Add buttons to quickly copy usernames and passwords.
- **Search and Filter**: Implement a search bar to filter the password list.
- **Master Password**: Secure the application with a master password.
- **Improved UI/UX**: Enhance visual feedback, loading states, and error notifications.

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
