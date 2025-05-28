# GopherNet 🦫

A modern platform for managing and monitoring gopher burrow rentals. GopherNet provides a robust API for managing burrow occupancy, tracking burrow statistics, and generating automated reports.

## Features

### Core Features
- 🏠 Burrow Management System
- 📊 Real-time Burrow Statistics
- 🔄 Automated Burrow Maintenance
- 📈 Periodic System Reports
- 🐳 Docker Support

### Bonus Features
- 💾 **Data Persistence**: All burrow data is automatically persisted between server restarts
- 📝 **Structured Logging**: Comprehensive logging using Zap logger with debug/production modes
- ⚙️ **Configurable Settings**: Flexible configuration via YAML with environment variable overrides
- 📊 **Enhanced Monitoring**: Detailed burrow statistics and automated reporting
- 🔄 **Smart Burrow Management**: Automatic depth updates and age-based cleanup
- 🧪 **Comprehensive Testing**: Extensive test coverage with mock-based testing
- 📚 **API Documentation**: Swagger/OpenAPI documentation for all endpoints

## Prerequisites

- Go 1.21 or later
- Make
- Docker (optional, for containerized deployment)

## Quick Start

### Local Development

1. Clone the repository:
```bash
git clone https://github.com/lebleuciel/GopherNet
cd gophernet
```

2. Run the setup command (this will handle everything):
```bash
make all
```

This single command will:
- Install all Go dependencies
- Install mockgen tool
- Generate mocks
- Run ent migrations
- Build and run the application

### Docker Deployment

1. Build the Docker image:
```bash
docker build -t gophernet .
```

2. Run the container:
```bash
docker run -p 8080:8080 gophernet
```

## Data Persistence

GopherNet automatically handles data persistence:

- On first run, the system loads initial burrow data from `data/initial.json`
- On subsequent runs, the system resumes the previous state from the database
- All burrow modifications (depth, occupancy, etc.) are persisted
- System reports are saved in the `reports` directory

## Logging

The application uses Zap logger with two modes:

- **Debug Mode**: Console-based logging with detailed information
- **Production Mode**: JSON-formatted logs with essential information

Configure logging in `config.yaml`:
```yaml
logger:
  debug: true  # Set to false for production mode
```

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

logger:
  debug: true
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

## Development

The project uses several tools that are automatically handled by the setup process:

- **ent**: For database schema and migrations
- **mockgen**: For generating mock interfaces
- **swag**: For Swagger documentation

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
