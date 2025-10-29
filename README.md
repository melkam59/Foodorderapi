# Foodorderapi

A Go-based RESTful API backend scaffold for a food ordering system. This project provides a starting point for building endpoints, middleware, and core domain logic for menus, orders, users, and other features you may add over time.

- Language: Go
- Default branch: `main`
- Repository: https://github.com/melkam59/Foodorderapi

> Note: This README describes the current repository layout and offers sensible defaults to run and extend the project. Update sections (features, config, API docs) as you implement specifics.

## Contents

```
.
├── .github/                 # GitHub configs/workflows
├── .gitignore
├── .idea/                   # IDE settings
├── docker-compose.yml       # Docker Compose services (app/db/etc.)
├── dokcerifle               # Dockerfile (currently misspelled)
├── go.mod                   # Go module definition
├── go.sum
├── internals/               # Core domain/application code
├── main/                    # Application entrypoint (main package)
├── middleware/              # HTTP middleware
├── qodana.yaml              # Qodana (JetBrains) static analysis config
├── routes/                  # HTTP route definitions
└── utils/                   # Shared helpers/utilities
```

## Getting Started

### Prerequisites

- Go (recommend Go 1.21+)
- Docker and Docker Compose (optional, for containerized runs)

### Clone

```bash
git clone https://github.com/melkam59/Foodorderapi.git
cd Foodorderapi
```

### Run locally (Go)

```bash
# Download dependencies
go mod download

# Run the server
go run ./main

# Or build a binary
go build -o bin/server ./main
./bin/server
```

### Run with Docker

This repo includes a Docker build file currently named `dokcerifle`. You can either:
- Build with the custom filename, or
- Rename it to `Dockerfile` before building.

Build directly using the existing filename:

```bash
docker build -f dokcerifle -t foodorderapi:dev .
docker run --rm -p 8080:8080 --env-file .env foodorderapi:dev
```

Or with Docker Compose:

```bash
# Adjust environment variables as needed (see Configuration)
docker compose up --build
# or
docker-compose up --build
```

## Configuration

Use environment variables to configure the service (create a `.env` file if you prefer). Common patterns:

```env
# Server
PORT=8080

# Database (choose one style your code supports)
DATABASE_URL=postgres://user:pass@localhost:5432/foodorder?sslmode=disable
# or granular:
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=foodorder
```

Adjust these to match your environment and any database or services defined in `docker-compose.yml`.

## API Endpoints

Base URL: http://localhost:8080

Notes:
- Routes marked as "protected" are behind token validation middleware (ValidateToken) and require an Authorization header (e.g., `Authorization: Bearer <token>`).
- Path params like `:id` and `:categoryid` are placeholders to be replaced with real values.

### Admin Authentication
- POST `/admin/` — Sign up an admin
- POST `/admin/login` — Admin login

### Admin (protected)
- PATCH `/admins/update/:id` — Update an admin by ID
- DELETE `/admins/delete/:id` — Delete an admin by ID
- POST `/admins/logout` — Logout

### Admin Managing Merchants (protected)
- POST `/admins/signupmerchant` — Create a merchant
- GET `/admins/allmerchant/` — List all merchants
- GET `/admins/singlemerchant/:id` — Get a merchant by ID
- PATCH `/admins/updatemerchant/:id` — Update a merchant by ID
- DELETE `/admins/deletemerchant/:id` — Delete a merchant by ID
- PUT `/admins/deactivatemerchant/:id` — Deactivate a merchant
- PUT `/admins/activatemerchant/:id` — Activate a merchant

### Merchant Authentication
- POST `/merchant/` — Merchant sign-in
- POST `/merchant/forgetpassword` — Merchant forgot password
- POST `/merchant/numberofmenubycategory` — Count menus per category

### Merchant Self-service (protected)
- POST `/merchants/logout` — Logout
- PATCH `/merchants/updateprofile` — Update merchant profile
- GET `/merchants/me` — Get current merchant profile

### Menu Management (merchant, protected)
- POST `/merchantmenu/addnewmenu` — Create a menu item
- GET `/merchantmenu/getallmenus` — List all menus (for merchant)
- GET `/merchantmenu/getsinglemenu/:id` — Get a menu by ID
- PATCH `/merchantmenu/updatemenu/:id` — Update a menu by ID
- DELETE `/merchantmenu/deletemenu/:id` — Delete a menu by ID

### Category Management (merchant, protected)
- POST `/category/new` — Create category
- GET `/category/all` — List categories
- PATCH `/category/:id` — Update category by ID
- DELETE `/category/:id` — Delete category by ID
- GET `/category/foods/:categoryid` — List foods in a category (merchant scope)
- GET `/category/numberofcategories` — Count categories (merchant scope)
- GET `/category/numberoffoods` — Count foods (merchant scope)

### User-facing Endpoints
- POST `/user/getmerchantdetail` — Get merchant by shortcode
- POST `/user/displayallmenu` — List menus
- POST `/user/displayallcategory` — List categories
- POST `/user/menubycategory/:categoryid` — List menus by category
- POST `/user/numberofmenubycategory` — Count menus per category
- POST `/user/fetchmenusbyfastingstatus` — List menus filtered by fasting status
- POST `/user/numberofcategories` — Count categories
- POST `/user/numberoffoods` — Count foods

## Development

- `main/`: Program entrypoint and bootstrap (HTTP server, wiring).
- `routes/`: Route registration and HTTP handlers mapping.
- `middleware/`: Cross-cutting HTTP middleware (auth, logging, CORS, recovery).
- `internals/`: Domain logic, services, repositories, and entities.
- `utils/`: Common helpers used across packages.

### Testing

```bash
go test ./...
```

### Linting and Static Analysis (Qodana)

The repository includes `qodana.yaml` for JetBrains Qodana. See the official docs to run locally via Docker:
- Qodana for Go: https://www.jetbrains.com/help/qodana/qodana-go.html

Example run (adjust paths as needed):

```bash
docker run --rm -it \
  -v "$PWD:/data/project" \
  -v "$PWD/.qodana:/data/results" \
  -p 8080:8080 \
  jetbrains/qodana-go:latest
```

## Troubleshooting

- Build errors related to the Docker file name:
  - Either rename `dokcerifle` to `Dockerfile` or pass `-f dokcerifle` to `docker build`.
- Port conflicts:
  - Change `PORT` in `.env` or your run command, and update port mappings in `docker-compose.yml`.

## Contributing

1. Fork the repo and create a feature branch.
2. Add tests for your changes when appropriate.
3. Run `go test ./...` and static checks before submitting.
4. Open a pull request with a clear description and context.

## License

No license has been specified in the repository. If you plan to share or reuse this code, consider adding a license (e.g., MIT, Apache-2.0).

## Acknowledgements

- Go community and standard library
- JetBrains Qodana for static analysis