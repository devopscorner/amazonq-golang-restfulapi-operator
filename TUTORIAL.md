# 🎓 RestAPI Operator Tutorial: From Zero to Production

## Tutorial Overview
This tutorial walks you through building and deploying a complete MVC+R pattern application using the RestAPI Operator.

## 🏗️ Tutorial Project: Guestbook Application

We'll build a guestbook application with:
- **Model**: User data management
- **View**: Web interface
- **Controller**: API endpoints
- **Repository**: Database operations

### Step 1: Environment Preparation

#### 1.1 Start Minikube (Local Development)
```bash
# Start local cluster
minikube start --memory=4096 --cpus=2

# Enable metrics server for auto-scaling
minikube addons enable metrics-server

# Verify cluster
kubectl cluster-info
```

#### 1.2 Deploy Required Dependencies
```bash
# Create namespace
kubectl create namespace guestbook

# Deploy PostgreSQL
kubectl apply -f - <<EOF
apiVersion: apps/v1
kind: Deployment
metadata:
  name: postgres
  namespace: guestbook
spec:
  replicas: 1
  selector:
    matchLabels:
      app: postgres
  template:
    metadata:
      labels:
        app: postgres
    spec:
      containers:
      - name: postgres
        image: postgres:13
        env:
        - name: POSTGRES_DB
          value: guestbook
        - name: POSTGRES_USER
          value: admin
        - name: POSTGRES_PASSWORD
          value: password
        ports:
        - containerPort: 5432
---
apiVersion: v1
kind: Service
metadata:
  name: postgres
  namespace: guestbook
spec:
  selector:
    app: postgres
  ports:
  - port: 5432
    targetPort: 5432
EOF
```

### Step 2: Deploy RestAPI Operator

#### 2.1 Install Operator
```bash
# Clone repository
git clone <your-repo>
cd amazonq-golang-restfulapi-operator

# Deploy operator
./deploy.sh

# Verify operator is running
kubectl get pods -n amazonq-golang-restfulapi-operator-system
```

#### 2.2 Verify Installation
```bash
# Check CRDs
kubectl get crd | grep restapi

# Check operator logs
kubectl logs -f deployment/amazonq-golang-restfulapi-operator-controller-manager -n amazonq-golang-restfulapi-operator-system
```

### Step 3: Create Guestbook Application

#### 3.1 Basic Configuration
Create `guestbook-basic.yaml`:
```yaml
apiVersion: apps.aws.com/v1
kind: RestAPI
metadata:
  name: guestbook
  namespace: guestbook
spec:
  image: "nginx:1.21"
  replicas: 1
  envVars:
    DATABASE_URL: "postgresql://admin:password@postgres:5432/guestbook"
    APP_ENV: "development"
  
  # MVC+R Components
  model:
    enabled: true
    image: "guestbook/model:v1.0.0"
    port: 8080
    envVars:
      COMPONENT_TYPE: "model"
      DB_POOL_SIZE: "5"
  
  view:
    enabled: true
    image: "guestbook/view:v1.0.0"
    port: 3000
    envVars:
      COMPONENT_TYPE: "view"
      API_BASE_URL: "http://guestbook-controller-svc:8081"
  
  controller:
    enabled: true
    image: "guestbook/controller:v1.0.0"
    port: 8081
    envVars:
      COMPONENT_TYPE: "controller"
      MODEL_SERVICE_URL: "http://guestbook-model-svc:8080"
      REPO_SERVICE_URL: "http://guestbook-repository-svc:8082"
  
  repository:
    enabled: true
    image: "guestbook/repository:v1.0.0"
    port: 8082
    envVars:
      COMPONENT_TYPE: "repository"
      DB_CONNECTION_TIMEOUT: "30"
```

#### 3.2 Deploy Basic Application
```bash
# Apply configuration
kubectl apply -f guestbook-basic.yaml

# Monitor deployment
kubectl get restapi -n guestbook -w

# Check created resources
kubectl get all -n guestbook -l app=guestbook
```

#### 3.3 Verify Deployment
```bash
# Check pods status
kubectl get pods -n guestbook

# Check services
kubectl get svc -n guestbook

# Test connectivity
kubectl port-forward -n guestbook svc/guestbook-view-svc 3000:3000
# Open browser: http://localhost:3000
```

### Step 4: Add Health Monitoring

#### 4.1 Update with Health Checks
Create `guestbook-health.yaml`:
```yaml
apiVersion: apps.aws.com/v1
kind: RestAPI
metadata:
  name: guestbook
  namespace: guestbook
spec:
  image: "nginx:1.21"
  replicas: 2
  envVars:
    DATABASE_URL: "postgresql://admin:password@postgres:5432/guestbook"
    APP_ENV: "development"
  
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
  
  # Add health monitoring
  healthCheck:
    enabled: true
    path: "/health"
    initialDelaySeconds: 30
    periodSeconds: 10
    timeoutSeconds: 5
    failureThreshold: 3
```

#### 4.2 Apply Health Configuration
```bash
# Update application
kubectl apply -f guestbook-health.yaml

# Check pod health
kubectl describe pod -n guestbook -l app=guestbook

# Monitor health checks
kubectl get events -n guestbook --sort-by=.metadata.creationTimestamp
```

### Step 5: Enable Auto-scaling

#### 5.1 Configure Auto-scaling
Create `guestbook-autoscale.yaml`:
```yaml
apiVersion: apps.aws.com/v1
kind: RestAPI
metadata:
  name: guestbook
  namespace: guestbook
spec:
  image: "nginx:1.21"
  replicas: 2
  envVars:
    DATABASE_URL: "postgresql://admin:password@postgres:5432/guestbook"
    APP_ENV: "production"
  
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
  
  # Health monitoring
  healthCheck:
    enabled: true
    path: "/health"
    initialDelaySeconds: 30
    periodSeconds: 10
    timeoutSeconds: 5
    failureThreshold: 3
  
  # Auto-scaling configuration
  autoScaling:
    enabled: true
    minReplicas: 2
    maxReplicas: 10
    targetCPUUtilization: 70
    targetMemoryUtilization: 80
```

#### 5.2 Test Auto-scaling
```bash
# Apply auto-scaling
kubectl apply -f guestbook-autoscale.yaml

# Check HPA status
kubectl get hpa -n guestbook

# Generate load to test scaling
kubectl run -i --tty load-generator --rm --image=busybox --restart=Never -- /bin/sh
# Inside pod:
while true; do wget -q -O- http://guestbook-controller-svc.guestbook:8081/api/entries; done

# Monitor scaling in another terminal
kubectl get hpa -n guestbook -w
kubectl get pods -n guestbook -w
```

### Step 6: Blue-Green Deployment

#### 6.1 Enable Blue-Green
Create `guestbook-bluegreen.yaml`:
```yaml
apiVersion: apps.aws.com/v1
kind: RestAPI
metadata:
  name: guestbook
  namespace: guestbook
spec:
  image: "nginx:1.21"
  replicas: 3
  envVars:
    DATABASE_URL: "postgresql://admin:password@postgres:5432/guestbook"
    APP_ENV: "production"
  
  model:
    enabled: true
    image: "guestbook/model:v1.1.0"  # Updated version
    port: 8080
  
  view:
    enabled: true
    image: "guestbook/view:v1.1.0"   # Updated version
    port: 3000
  
  controller:
    enabled: true
    image: "guestbook/controller:v1.1.0"  # Updated version
    port: 8081
  
  repository:
    enabled: true
    image: "guestbook/repository:v1.1.0"  # Updated version
    port: 8082
  
  healthCheck:
    enabled: true
    path: "/health"
    initialDelaySeconds: 30
    periodSeconds: 10
    timeoutSeconds: 5
    failureThreshold: 3
  
  autoScaling:
    enabled: true
    minReplicas: 2
    maxReplicas: 10
    targetCPUUtilization: 70
  
  # Blue-green deployment
  blueGreen:
    enabled: true
    strategy: "automatic"
    promotionTimeout: 300
```

#### 6.2 Deploy with Blue-Green
```bash
# Apply blue-green configuration
kubectl apply -f guestbook-bluegreen.yaml

# Monitor deployment
kubectl get restapi guestbook -n guestbook -o jsonpath='{.status.activeEnvironment}'

# Check both environments
kubectl get deployments -n guestbook -l app=guestbook

# Monitor traffic switching
kubectl get svc -n guestbook -l app=guestbook -o wide
```

### Step 7: Production Monitoring

#### 7.1 Set Up Monitoring
```bash
# Check application status
kubectl get restapi guestbook -n guestbook -o yaml

# Monitor all components
kubectl get pods,svc,hpa,deployments -n guestbook -l app=guestbook

# Check logs
kubectl logs -f deployment/guestbook-model -n guestbook
kubectl logs -f deployment/guestbook-controller -n guestbook
```

#### 7.2 Performance Testing
```bash
# Install hey for load testing
go install github.com/rakyll/hey@latest

# Port forward to test
kubectl port-forward -n guestbook svc/guestbook-controller-svc 8081:8081 &

# Run load test
hey -n 1000 -c 10 http://localhost:8081/api/entries

# Monitor scaling response
kubectl get hpa -n guestbook -w
```

### Step 8: Troubleshooting Exercise

#### 8.1 Simulate Failures
```bash
# Scale down database to simulate failure
kubectl scale deployment postgres -n guestbook --replicas=0

# Check application response
kubectl get pods -n guestbook
kubectl describe restapi guestbook -n guestbook

# Restore database
kubectl scale deployment postgres -n guestbook --replicas=1
```

#### 8.2 Debug Common Issues
```bash
# Check operator logs
kubectl logs -f deployment/amazonq-golang-restfulapi-operator-controller-manager -n amazonq-golang-restfulapi-operator-system

# Check RestAPI events
kubectl describe restapi guestbook -n guestbook

# Check pod logs
kubectl logs -f deployment/guestbook-model -n guestbook
```

## 🎯 Tutorial Completion Checklist

- [ ] Operator deployed and running
- [ ] Basic RestAPI application deployed
- [ ] Health checks configured and working
- [ ] Auto-scaling enabled and tested
- [ ] Blue-green deployment configured
- [ ] Monitoring and logging set up
- [ ] Load testing completed
- [ ] Troubleshooting scenarios practiced

## 🚀 Next Steps

1. **Custom Images**: Build your own MVC+R component images
2. **CI/CD Integration**: Integrate with GitOps workflows
3. **Multi-Environment**: Set up staging and production environments
4. **Advanced Monitoring**: Add Prometheus and Grafana
5. **Security**: Implement RBAC and network policies

## 📚 Additional Resources

- [Kubernetes Documentation](https://kubernetes.io/docs/)
- [Operator SDK Guide](https://sdk.operatorframework.io/)
- [MVC Pattern Best Practices](https://martinfowler.com/eaaCatalog/modelViewController.html)
- [Blue-Green Deployment Strategies](https://martinfowler.com/bliki/BlueGreenDeployment.html)