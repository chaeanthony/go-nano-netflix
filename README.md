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

## Getting Started

### Prerequisites

To run this API, you need:

- Go (this project was built using 1.23.2)
- PostgreSQL

### Installation

1. Clone the repository:

```bash
git clone https://github.com/chaeanthony/go-netflix.git
cd go-netflix
```

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

### User Management

### Title Management
