# HOW-TO: RestAPI Operator for MVC+R Pattern

## 🚀 Quick Start Guide

### Prerequisites
```bash
# Required tools
kubectl version --client  # v1.20+
docker --version          # 17.03+
go version                # 1.21+ (for development)
```

### Step 1: Deploy the Operator
```bash
# Clone and navigate
git clone <repository>
cd amazonq-golang-restfulapi-operator

# Deploy operator
./deploy.sh
```

### Step 2: Create Your First RestAPI
```bash
# Apply sample
kubectl apply -f config/samples/apps_v1_restapi.yaml

# Check status
kubectl get restapi
kubectl get pods
```

## 📋 Complete Step-by-Step Tutorial

### Phase 1: Environment Setup

#### 1.1 Verify Kubernetes Cluster
```bash
# Check cluster connection
kubectl cluster-info

# Verify permissions
kubectl auth can-i create customresourcedefinitions
```

#### 1.2 Install Operator
```bash
# Install CRDs
make install

# Deploy operator
make deploy

# Verify operator is running
kubectl get pods -n amazonq-golang-restfulapi-operator-system
```

### Phase 2: Basic RestAPI Deployment

#### 2.1 Simple MVC+R Application
Create `my-restapi.yaml`:
```yaml
apiVersion: apps.aws.com/v1
kind: RestAPI
metadata:
  name: my-app
spec:
  image: "nginx:1.21"
  replicas: 2
  
  model:
    enabled: true
    image: "myapp/model:v1"
    port: 8080
  
  view:
    enabled: true
    image: "myapp/view:v1"
    port: 3000
  
  controller:
    enabled: true
    image: "myapp/controller:v1"
    port: 8081
  
  repository:
    enabled: true
    image: "myapp/repository:v1"
    port: 8082
```

#### 2.2 Deploy and Monitor
```bash
# Deploy
kubectl apply -f my-restapi.yaml

# Monitor deployment
kubectl get restapi my-app -w

# Check created resources
kubectl get deployments,services,pods -l app=my-app
```

### Phase 3: Advanced Features

#### 3.1 Enable Auto-scaling
```yaml
spec:
  autoScaling:
    enabled: true
    minReplicas: 2
    maxReplicas: 10
    targetCPUUtilization: 70
    targetMemoryUtilization: 80
```

#### 3.2 Configure Health Monitoring
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

#### 3.3 Enable Blue-Green Deployment
```yaml
spec:
  blueGreen:
    enabled: true
    strategy: "automatic"
    promotionTimeout: 300
```

### Phase 4: Production Configuration

#### 4.1 Environment Variables
```yaml
spec:
  envVars:
    DATABASE_URL: "postgresql://db:5432/myapp"
    REDIS_URL: "redis://cache:6379"
    LOG_LEVEL: "info"
  
  model:
    enabled: true
    envVars:
      DB_POOL_SIZE: "20"
      CACHE_TTL: "300"
```

#### 4.2 Resource Management
```yaml
spec:
  model:
    enabled: true
    image: "myapp/model:v2.0.0"
    port: 8080
    resources:
      requests:
        memory: "256Mi"
        cpu: "250m"
      limits:
        memory: "512Mi"
        cpu: "500m"
```

## 🔧 Common Operations

### Scaling Components
```bash
# Scale specific component
kubectl patch restapi my-app --type='merge' -p='{"spec":{"replicas":5}}'

# Enable/disable components
kubectl patch restapi my-app --type='merge' -p='{"spec":{"view":{"enabled":false}}}'
```

### Monitoring and Debugging
```bash
# Check RestAPI status
kubectl describe restapi my-app

# View operator logs
kubectl logs -f deployment/amazonq-golang-restfulapi-operator-controller-manager -n amazonq-golang-restfulapi-operator-system

# Check component health
kubectl get pods -l app=my-app
kubectl describe pod <pod-name>
```

### Blue-Green Deployment Management
```bash
# Check current environment
kubectl get restapi my-app -o jsonpath='{.status.activeEnvironment}'

# Manual promotion (if strategy is manual)
kubectl patch restapi my-app --type='merge' -p='{"spec":{"blueGreen":{"promote":true}}}'
```

## 📊 Monitoring and Observability

### Health Checks
```bash
# Check all components health
kubectl get pods -l app=my-app -o wide

# Test health endpoints
kubectl port-forward svc/my-app-model-svc 8080:8080
curl http://localhost:8080/health
```

### Auto-scaling Status
```bash
# Check HPA status
kubectl get hpa -l app=my-app

# View scaling events
kubectl describe hpa my-app-model-hpa
```

### Service Discovery
```bash
# List all services
kubectl get svc -l app=my-app

# Test service connectivity
kubectl run test-pod --image=curlimages/curl -it --rm -- sh
# Inside pod: curl http://my-app-model-svc:8080/health
```

## 🛠️ Troubleshooting Guide

### Common Issues

#### Operator Not Starting
```bash
# Check operator logs
kubectl logs -f deployment/amazonq-golang-restfulapi-operator-controller-manager -n amazonq-golang-restfulapi-operator-system

# Verify CRDs installed
kubectl get crd restapis.apps.aws.com
```

#### Components Not Deploying
```bash
# Check RestAPI status
kubectl describe restapi my-app

# Verify images exist
docker pull myapp/model:v1

# Check resource constraints
kubectl describe nodes
```

#### Health Checks Failing
```bash
# Check probe configuration
kubectl describe pod <pod-name>

# Test health endpoint manually
kubectl exec -it <pod-name> -- curl localhost:8080/health
```

### Performance Tuning

#### Optimize Auto-scaling
```yaml
spec:
  autoScaling:
    enabled: true
    minReplicas: 3          # Higher minimum for stability
    maxReplicas: 20         # Higher maximum for peak loads
    targetCPUUtilization: 60 # Lower threshold for faster scaling
```

#### Resource Allocation
```yaml
spec:
  model:
    resources:
      requests:
        memory: "512Mi"     # Higher requests for guaranteed resources
        cpu: "500m"
      limits:
        memory: "1Gi"       # Reasonable limits
        cpu: "1000m"
```

## 🔄 Upgrade and Maintenance

### Updating Components
```bash
# Update component image
kubectl patch restapi my-app --type='merge' -p='{"spec":{"model":{"image":"myapp/model:v2.0.0"}}}'

# Rolling update with blue-green
kubectl patch restapi my-app --type='merge' -p='{"spec":{"blueGreen":{"enabled":true}}}'
```

### Backup and Recovery
```bash
# Backup RestAPI configuration
kubectl get restapi my-app -o yaml > my-app-backup.yaml

# Restore from backup
kubectl apply -f my-app-backup.yaml
```

### Cleanup
```bash
# Delete RestAPI (cascades to all components)
kubectl delete restapi my-app

# Uninstall operator
make undeploy
make uninstall
```

## 📚 Advanced Examples

### Multi-Environment Setup
```yaml
# Production
apiVersion: apps.aws.com/v1
kind: RestAPI
metadata:
  name: my-app-prod
  namespace: production
spec:
  replicas: 5
  autoScaling:
    enabled: true
    minReplicas: 5
    maxReplicas: 50
  blueGreen:
    enabled: true
    strategy: "manual"

---
# Staging
apiVersion: apps.aws.com/v1
kind: RestAPI
metadata:
  name: my-app-staging
  namespace: staging
spec:
  replicas: 2
  autoScaling:
    enabled: true
    minReplicas: 1
    maxReplicas: 10
```

### Microservices Architecture
```yaml
# User Service
apiVersion: apps.aws.com/v1
kind: RestAPI
metadata:
  name: user-service
spec:
  model:
    enabled: true
    image: "microservices/user-model:v1"
  repository:
    enabled: true
    image: "microservices/user-repo:v1"
  # Disable view and controller for API-only service
  view:
    enabled: false
  controller:
    enabled: false
```

## 🎯 Best Practices

1. **Resource Planning**: Always set resource requests/limits
2. **Health Checks**: Implement proper health endpoints
3. **Environment Variables**: Use ConfigMaps/Secrets for sensitive data
4. **Monitoring**: Enable metrics and logging
5. **Testing**: Test blue-green deployments in staging first
6. **Backup**: Regular backup of RestAPI configurations
7. **Security**: Use least-privilege RBAC policies

## 📞 Support

- **Logs**: `kubectl logs -f deployment/amazonq-golang-restfulapi-operator-controller-manager -n amazonq-golang-restfulapi-operator-system`
- **Status**: `kubectl describe restapi <name>`
- **Events**: `kubectl get events --sort-by=.metadata.creationTimestamp`