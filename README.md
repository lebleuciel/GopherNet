# GopherNet 🦫

A modern platform for managing and monitoring gopher burrow rentals. GopherNet provides a robust API for managing burrow occupancy, tracking burrow statistics, and generating automated reports.

## Features

- 🏠 Burrow Management System
- 📊 Real-time Burrow Statistics
- 🔄 Automated Burrow Maintenance
- 📈 Periodic System Reports
- 🐳 Docker Support
- 🧪 Comprehensive Test Coverage
- 📚 Swagger API Documentation
- 💾 Data Persistence Between Runs

## Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose
- PostgreSQL 14 or higher

## Quick Start

### Using Docker

1. Clone the repository:
```bash
git clone https://github.com/lebleuciel/GopherNet
cd gophernet
```

2. Start the services:
```bash
make docker-up
```

3. The API will be available at `http://localhost:8080`

### Manual Setup

1. Install dependencies:
```bash
make deps
```

2. Set up the database:
```bash
make db-create
```

3. Build and run:
```bash
make build
make run
```

## Data Persistence

GopherNet automatically handles data persistence:

- On first run, the system loads initial burrow data from `data/initial.json`
- On subsequent runs, the system resumes the previous state from the database
- All burrow modifications (depth, occupancy, etc.) are persisted
- System reports are saved in the `reports` directory

## API Endpoints

### Get Gopher Status
```bash
curl -X GET http://localhost:8080/api/v1/gopher
```

### Get All Burrows Status
```bash
curl -X GET http://localhost:8080/api/v1/burrows/status
```

### Rent a Burrow
```bash
curl -X POST http://localhost:8080/api/v1/burrows/1/rent
```

### Release a Burrow
```bash
curl -X POST http://localhost:8080/api/v1/burrows/1/release
```

## Docker Commands

- Build images: `make docker-build`
- Start services: `make docker-up`
- View logs: `make docker-logs`
- Stop services: `make docker-down`
- Clean up: `make docker-clean`

## Testing

Run the test suite:
```bash
make test
```

Generate and update mocks:
```bash
make install-mockgen
make generate-mocks
```

## Configuration

The application uses a `config.yaml` file for configuration. Here's the default configuration:

```yaml
database:
  host: db
  port: 5432
  user: postgres
  password: postgres
  database: gophernet

scheduler:
  report_interval: 2m
  update_interval: 1m
  max_burrow_age: 1440
  depth_increment: 0.009
```

## Initial Data

The system comes with a set of initial burrows. Here's the sample `initial.json`:

```json
[
  {
    "name": "The Underground Palace",
    "depth": 2.5,
    "width": 1.2,
    "occupied": true,
    "age": 10
  },
  {
    "name": "Tunnel of Mystery",
    "depth": 1.8,
    "width": 1.1,
    "occupied": false,
    "age": 30
  },
  {
    "name": "The Molehole",
    "depth": 3.0,
    "width": 1.3,
    "occupied": true,
    "age": 50
  },
  {
    "name": "The Deep Den",
    "depth": 2.2,
    "width": 1.2,
    "occupied": false,
    "age": 40
  },
  {
    "name": "Surface Level Statistics",
    "depth": 0,
    "width": 1.3,
    "occupied": true,
    "age": 5
  }
]
```

## Project Structure

```
.
├── cmd/            # Application entry points
├── pkg/            # Core packages
│   ├── app/        # Business logic
│   ├── config/     # Configuration
│   ├── controller/ # HTTP controllers
│   ├── db/         # Database models and migrations
│   ├── dto/        # Data transfer objects
│   ├── mocks/      # Generated mocks
│   ├── models/     # Domain models
│   ├── repo/       # Repository interfaces
│   └── utils/      # Utility functions
├── data/           # Data files
├── docs/           # Documentation
└── server/         # HTTP server setup
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
