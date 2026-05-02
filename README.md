# Beacon

Beacon is a simple tool to watch for Docker deployment notification in CI/CD workflow - designed for machines under controlled networks.

```mermaid
flowchart LR
    A[Runner] -->|Publishes message| B[GCP Pub/Sub Topic]

    subgraph PN[Controlled Network]
        C[Beacon]
        D[Docker]
        C -->|Trigger Deployment| D
    end

    B -->|Subscription / Change Event| C
```

## Setup

### Prerequisites

**1. GCP Authentication**

Beacon uses [Application Default Credentials (ADC)](https://cloud.google.com/docs/authentication/application-default-credentials) to authenticate with GCP Pub/Sub. Authenticate before running Beacon:

```bash
gcloud auth application-default login
```

Alternatively, set the `GOOGLE_APPLICATION_CREDENTIALS` environment variable to the path of a service account key file:
```bash
export GOOGLE_APPLICATION_CREDENTIALS="/path/to/service-account-key.json"
```

**2. Create a configuration file**

Create a `config.yaml` file. See the [Configuration](#configuration) section below for all available options.

**3. Build the binary (if running from source)**

```bash
go build -o beacon .
```

---

### VM-hosted

Grant Docker access to the user running Beacon, then run the binary with your config file.

- Option 1: Add the current user to the `docker` group
```bash
sudo usermod -aG docker $USER
newgrp docker

./beacon -config config.yaml
```

- Option 2: Run as root
```bash
sudo ./beacon -config config.yaml
```

---

### Docker

Beacon must have access to the host's Docker socket (Docker-outside-of-Docker). Mount the socket and your config file into the container:

```bash
docker run \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v $HOME/.docker:/root/.docker:ro \
  -v ./config.yaml:/app/config.yaml:ro \
  beacon -config /app/config.yaml
```

## Configuration
```yaml
# One GCP project is supported per app instance
gcp-project-id: "your-project-id"

# Define multiple consumers
consumers:
  my-topic-consumer: # Any ID you want
    pubsub-subscription-id: "your-subscription-id"

    deduplication:
      enabled: false
      
      # When enabled, deployment messages within a 5-minute window results in a single deployment trigger
      time-window: "5m"
    
    trigger-commands:
      - 'echo "Triggering Docker deployment..."'
      - 'docker stack deploy --with-registry-auth -c docker-compose.yml myapp'
```