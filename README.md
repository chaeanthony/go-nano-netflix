# Mini-netflix clone

Backend for a mini netflix clone. Use as starter code for similar applications. Built using Go (Golang), PostgreSQL.

Features:

- Follows RESTful principles.
- JWT Access Tokens
- Refresh Tokens

## Table of Contents

- [Getting Started](#getting-started)
- [API](#api)
  - [User Management](#user-management)
  - [Title Management](#title-management)
  - [Watchlist Management](#watchlist-management)

## Getting Started

### Prerequisites

To run this API, you need:

- Go (this project was built using 1.23.2)
- PostgreSQL

### Installation

1. Fork/Clone the repository:

2. Install dependencies:

- [Go](https://golang.org/doc/install)

```bash
go mod download
```

3. Create .env:

Copy the `.env.example` file to `.env` and fill in the values.

```bash
cp .env.example .env
```

Can use the following to generate a random 64-byte string encoded in Base64 as jwt secret

```bash
openssl rand -base64 64
```

4. Run server:

```bash
go run .
```

OR

```bash
go build -o out && ./out
```

## API

#### Base URL

The base URL for the API is `/api`. Use `/auth` for authentication APIs.

### Health Check

- **Path**: `/api/healthz`
- **Method**: `GET`
- **Description**: Returns API ready.

### User Management

#### Create User

- **Path**: `/auth/signup`
- **Method**: `POST`
- **Parameters**: `{"email": "test@email.com", "password": "123456"}`
- **Description**: Creates a new user.

#### User Login

- **Path**: `/auth/login`
- **Method**: `POST`
- **Parameters**: `{"email": "test@email.com", "password": "123456"}`
- **Description**: Authenticates a user and returns a jwt token and refresh token.

#### Refresh Token

- **Path**: `/auth/refresh`
- **Method**: `POST`
- **Description**: Accepts Bearer token in authorization header (this is a refresh token). Uses refresh token to create a new jwt token.

#### Revoke Token

- **Path**: `/auth/revoke`
- **Method**: `POST`
- **Description**: Revokes the user's refresh token.

### Title Management

#### Get All Titles

- **Path**: `/api/titles`
- **Method**: `GET`
- **Description**: Retrives all titles

#### Get Title by ID

- **Path**: `/api/titles/{titleId}`
- **Method**: `GET`
- **Description**: Retrieves details of a specific title based on its ID.

#### Get All Shows

- **Path**: `/api/shows`
- **Method**: `GET`
- **Description**: Retrieves all titles of type "show" available in the system.

#### Get All Movies

- **Path**: `/api/movies`
- **Method**: `GET`
- **Description**: Retrieves all titles of type "movie" available in the system.

### Watchlist Management

#### Get User's Watchlist

- **Path**: `/api/watchlist`
- **Headers**: `Authorization: Bearer <token>`
- **Method**: `GET`
- **Description**: Retrieves the watchlist for the authenticated user. Requires a valid bearer token for authentication.

#### Add Item to User's Watchlist

- **Path**: `/api/watchlist`
- **Headers**: `Authorization: Bearer <token>`
- **Method**: `POST`
- **Parameters**: `{"title_id": <int>}`
- **Description**: Retrieves the watchlist for the authenticated user. Requires a valid bearer token for authentication.
