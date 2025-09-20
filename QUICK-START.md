# ⚡ Quick Start Guide - RestAPI Operator

## 🎯 5-Minute Setup

### Prerequisites Check
```bash
kubectl version --client    # ✅ v1.20+
docker --version           # ✅ 17.03+
```

### 1. Deploy Operator (2 minutes)
```bash
# Clone and deploy
git clone <repository>
cd amazonq-golang-restfulapi-operator
./deploy.sh
```

### 2. Deploy Sample App (1 minute)
```bash
# Apply guestbook sample
kubectl apply -f config/samples/apps_v1_restapi.yaml

# Check status
kubectl get restapi
```

### 3. Verify Deployment (2 minutes)
```bash
# Check all resources
kubectl get pods,svc,hpa -l app=guestbook-restapi

# Test connectivity
kubectl port-forward svc/guestbook-restapi-view-svc 3000:3000
# Open: http://localhost:3000
```

## 🔧 Common Commands

### Check Status
```bash
kubectl get restapi -A
kubectl describe restapi <name>
```

### Scale Application
```bash
kubectl patch restapi <name> --type='merge' -p='{"spec":{"replicas":5}}'
```

### View Logs
```bash
kubectl logs -f deployment/<name>-model
kubectl logs -f deployment/amazonq-golang-restfulapi-operator-controller-manager -n amazonq-golang-restfulapi-operator-system
```

### Enable Features
```bash
# Enable auto-scaling
kubectl patch restapi <name> --type='merge' -p='{"spec":{"autoScaling":{"enabled":true,"maxReplicas":10}}}'

# Enable blue-green
kubectl patch restapi <name> --type='merge' -p='{"spec":{"blueGreen":{"enabled":true}}}'
```

## 🚨 Troubleshooting

### Operator Not Starting
```bash
kubectl logs -f deployment/amazonq-golang-restfulapi-operator-controller-manager -n amazonq-golang-restfulapi-operator-system
```

### Pods Not Running
```bash
kubectl describe restapi <name>
kubectl get events --sort-by=.metadata.creationTimestamp
```

### Health Checks Failing
```bash
kubectl describe pod <pod-name>
kubectl exec -it <pod-name> -- curl localhost:8080/health
```

## 📖 Next Steps

- Read [HOW-TO-USE.md](HOW-TO-USE.md) for detailed usage
- Follow [TUTORIAL.md](TUTORIAL.md) for hands-on learning
- Check [examples/](examples/) for more configurations