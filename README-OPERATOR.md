# RestAPI Operator - MVC+R Pattern Kubernetes Operator

A Kubernetes operator that manages MVC+R (Model, View, Controller, Repository) pattern RestAPI deployments with auto-scaling, health monitoring, and blue-green deployments.

## Features

- **MVC+R Pattern Management**: Deploy and manage Model, View, Controller, and Repository components separately
- **Auto-scaling**: Horizontal Pod Autoscaler (HPA) integration with CPU and memory metrics
- **Health Monitoring**: Configurable liveness and readiness probes
- **Blue-Green Deployments**: Zero-downtime deployments with traffic switching
- **Service Discovery**: Automatic service creation for each component

## Architecture

The operator manages four main components:

1. **Model**: Data layer handling business logic and data structures
2. **View**: Presentation layer for API responses and UI rendering
3. **Controller**: Request handling and routing logic
4. **Repository**: Data access layer for database operations

## Installation

### Prerequisites

- Kubernetes cluster (v1.20+)
- kubectl configured
- operator-sdk (optional, for development)

### Deploy the Operator

```bash
# Apply CRDs
make install

# Deploy the operator
make deploy

# Or build and deploy locally
make docker-build docker-push IMG=<your-registry>/restapi-operator:tag
make deploy IMG=<your-registry>/restapi-operator:tag
```

## Usage

### Basic RestAPI Resource

```yaml
apiVersion: apps.aws.com/v1
kind: RestAPI
metadata:
  name: guestbook-api
spec:
  image: "nginx:1.21"
  replicas: 2
  
  # MVC+R Components
  model:
    enabled: true
    image: "guestbook/model:v1.0.0"
    port: 8080
  
  view:
    enabled: true
    image: "guestbook/view:v1.0.0"
    port: 3000
  
  controller:
    enabled: true
    image: "guestbook/controller:v1.0.0"
    port: 8081
  
  repository:
    enabled: true
    image: "guestbook/repository:v1.0.0"
    port: 8082
```

### With Auto-scaling

```yaml
spec:
  autoScaling:
    enabled: true
    minReplicas: 1
    maxReplicas: 10
    targetCPUUtilization: 70
    targetMemoryUtilization: 80
```

### With Health Monitoring

```yaml
spec:
  healthCheck:
    enabled: true
    path: "/health"
    initialDelaySeconds: 30
    periodSeconds: 10
    timeoutSeconds: 5
    failureThreshold: 3
```

### With Blue-Green Deployment

```yaml
spec:
  blueGreen:
    enabled: true
    strategy: "automatic"
    promotionTimeout: 300
```

## Component Configuration

Each MVC+R component supports:

- **enabled**: Enable/disable the component
- **image**: Container image for the component
- **port**: Service port for the component
- **envVars**: Environment variables specific to the component

## Monitoring

The operator creates the following resources for each enabled component:

- **Deployment**: Manages pod replicas
- **Service**: Provides service discovery
- **HorizontalPodAutoscaler**: Manages auto-scaling (if enabled)

## Blue-Green Deployments

When blue-green deployment is enabled:

1. Both blue and green environments are maintained
2. Traffic is routed to the active environment
3. New deployments go to the inactive environment
4. Traffic switches after successful health checks
5. Old environment is kept for rollback capability

## Development

### Prerequisites

- Go 1.21+
- Docker
- kubectl
- operator-sdk

### Build and Test

```bash
# Generate manifests
make manifests

# Run tests
make test

# Build operator
make build

# Run locally (requires kubeconfig)
make run
```

### Generate CRDs

```bash
make manifests
```

## API Reference

### RestAPISpec

| Field | Type | Description |
|-------|------|-------------|
| image | string | Default container image |
| replicas | *int32 | Number of replicas |
| envVars | map[string]string | Global environment variables |
| model | ComponentSpec | Model component configuration |
| view | ComponentSpec | View component configuration |
| controller | ComponentSpec | Controller component configuration |
| repository | ComponentSpec | Repository component configuration |
| autoScaling | *AutoScalingSpec | Auto-scaling configuration |
| healthCheck | *HealthCheckSpec | Health check configuration |
| blueGreen | *BlueGreenSpec | Blue-green deployment configuration |

### ComponentSpec

| Field | Type | Description |
|-------|------|-------------|
| enabled | bool | Enable/disable component |
| image | string | Component container image |
| port | int32 | Service port |
| envVars | map[string]string | Component environment variables |

### AutoScalingSpec

| Field | Type | Description |
|-------|------|-------------|
| enabled | bool | Enable auto-scaling |
| minReplicas | *int32 | Minimum replicas |
| maxReplicas | int32 | Maximum replicas |
| targetCPUUtilization | *int32 | Target CPU utilization percentage |
| targetMemoryUtilization | *int32 | Target memory utilization percentage |

## Examples

See `config/samples/` for complete examples.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make changes and add tests
4. Run `make test` and `make manifests`
5. Submit a pull request

## License

Licensed under the Apache License, Version 2.0.