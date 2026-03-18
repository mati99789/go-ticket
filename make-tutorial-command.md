# Make Commands Guide

This document describes the available Make commands for managing the infrastructure and Go application.

## Overview

The following Make targets are available to streamline your development workflow. They handle both local development with manual infrastructure and containerized environments.

---

## Local Development Commands

### `make infra`

Starts the infrastructure stack required for the application.

**Use when:** You want to spin up databases, message brokers, and other dependent services locally.

```bash
make infra
```

---

### `make run`

Launches the Go application in local development mode.

**Use when:** You're developing features and want to run the application directly on your machine.

**Prerequisites:** Infrastructure must be running (`make infra`)

```bash
make run
```

---

### `make infra-down`

Stops and tears down the infrastructure stack.

**Use when:** You're finished developing and want to clean up local services.

```bash
make infra-down
```

---

## Monitoring & Debugging

### `make logs`

Displays real-time logs from the application and services.

**Use when:** You need to monitor application behavior, debug issues, or track request flow.

```bash
make logs
```

---

### `make logs-kafka`

Shows real-time logs specifically from the Kafka message broker.

**Use when:** You're debugging message publishing/consumption or investigating queue issues.

```bash
make logs-kafka
```

---

## Docker Commands

### `make docker`

Builds and starts the complete application stack in Docker containers.

**Use when:**

- You want an isolated, containerized environment
- You're testing production-like conditions locally
- You want everything to run with a single command

```bash
make docker
```

---

### `make docker-down`

Stops and removes all Docker containers and networks.

**Use when:** You're done working and want to clean up containerized resources.

```bash
make docker-down
```

---

## Typical Workflows

### Local Development

```bash
make infra          # Start services
make run            # Run application
make logs           # Monitor output
# ... develop ...
make infra-down     # Cleanup
```

### Docker Development

```bash
make docker         # Start everything
make logs           # Watch output
# ... develop ...
make docker-down    # Cleanup
```

### Debugging Kafka

```bash
make infra          # Start infrastructure
make logs-kafka     # Watch Kafka logs
# ... investigate ...
make infra-down     # Cleanup
```

---

## Notes

- **Ports & Configuration:** Ensure required ports are available before running `make infra` or `make docker`
- **Dependencies:** The Go application depends on services started by `make infra`
- **Data Persistence:** Review your configuration to understand if data persists between `down` and subsequent `up` commands
- **Multiple Terminals:** You may want to run `make logs` in a separate terminal while keeping `make run` or `make docker` in another
