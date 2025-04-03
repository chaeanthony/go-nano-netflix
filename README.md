# Mini-netflix clone

Mini netflix clone. Use as starter code for similar applications. Built using Go (Golang), PostgreSQL, React, Tailwind.

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
git clone https://github.com/chaeanthony/chirpy.git
cd chirpy
```

2. Install dependencies:

```bash
go mod tidy
```

3. Create .env:

```
DB_URL = "postgres database url"
PLATFORM = "dev" (prevent dangerous endpoints from being accessed in production)
JWT_SECRET = "jwt secret"
```

Can use the following to generate a random 64-byte string encoded in Base64 as jwt secret

```bash
openssl rand -base64 64
```

## API

#### Base URL

The base URL for the API is `/api`. Use `/auth` for authentication APIs.

### User Management

### Title Management
