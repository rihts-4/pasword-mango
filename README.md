# Pasword Mango

Pasword Mango is a secure and user-friendly password manager designed to help you manage your credentials effortlessly. It leverages a modern tech stack with a Go backend for performance, Clerk for robust authentication, Convex for real-time database capabilities, and a C++ Qt frontend for a native cross-platform experience.

## ✨ Features

- **Secure Credential Storage**: Safely store and organize your usernames, passwords, and other sensitive information.
- **Robust Authentication**: User management is handled securely by [Clerk](https://clerk.com/).
- **Real-time Database**: Utilizes [Convex](https://www.convex.dev/) as a real-time database to keep your data synced across devices.
- **Cross-Platform Desktop App**: The frontend is built with C++ and the [Qt framework](https://www.qt.io/), allowing it to run on Windows, macOS, and Linux.
- **High-Performance Backend**: The backend is written in [Go](https://go.dev/), ensuring speed and reliability.

## 🛠️ Tech Stack

- **Backend**: Go
- **Frontend**: C++ with Qt
- **Authentication**: Clerk
- **Database**: Convex

## 🚀 Getting Started

Follow these instructions to get a local copy up and running for development and testing purposes.

### Prerequisites

- [Go](https://go.dev/doc/install) (version 1.18 or newer)
- [Qt](https://www.qt.io/download) and a C++ compiler (e.g., GCC, Clang, MSVC)
- A [Clerk](https://clerk.com/) account for authentication keys.
- A [Convex](https://www.convex.dev/) account for the database.

### Installation & Setup

1.  **Clone the repository:**

    ```sh
    git clone https://github.com/your-username/pasword-mango.git
    cd pasword-mango
    ```

2.  **Backend (Go):**

    - Navigate to the backend directory.
    - Set up your environment variables for Clerk and Convex in a `.env` file.
      ```env
      # .env
      CLERK_SECRET_KEY="your_clerk_secret_key"
      CONVEX_DEPLOYMENT_URL="your_convex_deployment_url"
      ```
    - Install dependencies and run the server:
      ```sh
      go mod tidy
      go run main.go
      ```

3.  **Frontend (C++ & Qt):**
    - Navigate to the `frontend` directory (you may need to create this).
    - Use Qt Creator or your preferred C++ IDE to open the `.pro` or `CMakeLists.txt` file.
    - Build and run the project. The application will connect to the local Go backend.

## 📂 Project Structure

The project is organized with a clear separation between the frontend and backend code.
