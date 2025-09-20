# 🤝 Contributing to RestAPI Operator

## Quick Start for Contributors

### Prerequisites
```bash
go version    # 1.21+
kubectl       # v1.20+
docker        # 17.03+
operator-sdk  # latest
```

### Development Setup
```bash
# Fork and clone
git clone https://github.com/devopscorner/amazonq-golang-restfulapi-operator
cd amazonq-golang-restfulapi-operator

# Install dependencies
go mod tidy

# Run tests
make test

# Build locally
make build
```

## 🛠️ Development Workflow

### 1. Create Feature Branch
```bash
git checkout -b feature/your-feature-name
```

### 2. Make Changes
```bash
# Edit code
vim api/v1/restapi_types.go
vim internal/controller/restapi_controller.go

# Generate manifests
make manifests

# Run tests
make test

# Build and verify
make build
```

### 3. Test Changes
```bash
# Run operator locally
make run

# In another terminal, test with sample
kubectl apply -f config/samples/apps_v1_restapi.yaml
```

### 4. Submit PR
```bash
git add .
git commit -m "feat: add new feature"
git push origin feature/your-feature-name
```

## 📋 Contribution Areas

### Core Features
- [ ] **New Component Types**: Add support for additional MVC+R components
- [ ] **Advanced Scaling**: Implement custom metrics scaling
- [ ] **Multi-cluster**: Support for cross-cluster deployments
- [ ] **Security**: Enhanced RBAC and security policies

### Integrations
- [ ] **Service Mesh**: Istio/Linkerd integration
- [ ] **Monitoring**: Prometheus/Grafana dashboards
- [ ] **CI/CD**: GitOps workflow templates
- [ ] **Storage**: Persistent volume management

### Documentation
- [ ] **API Reference**: Complete API documentation
- [ ] **Examples**: More real-world examples
- [ ] **Tutorials**: Advanced use case tutorials
- [ ] **Troubleshooting**: Common issue solutions

## 🧪 Testing Guidelines

### Unit Tests
```bash
# Run all tests
make test

# Test specific package
go test ./internal/controller/...

# Test with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Integration Tests
```bash
# Run e2e tests
make test-e2e

# Test with local cluster
kind create cluster
make deploy
kubectl apply -f config/samples/
```

### Manual Testing
```bash
# Test basic functionality
kubectl apply -f examples/complete-example.yaml
kubectl get restapi -A
kubectl get pods,svc,hpa

# Test scaling
kubectl patch restapi guestbook-api --type='merge' -p='{"spec":{"replicas":5}}'

# Test blue-green
kubectl patch restapi guestbook-api --type='merge' -p='{"spec":{"blueGreen":{"enabled":true}}}'
```

## 📝 Code Standards

### Go Code Style
```go
// Use meaningful names
type RestAPIReconciler struct {
    client.Client
    Scheme *runtime.Scheme
}

// Add comments for exported functions
// ReconcileMVCRComponents manages MVC+R pattern deployments
func (r *RestAPIReconciler) ReconcileMVCRComponents(ctx context.Context, restAPI *restapiv1.RestAPI) error {
    // Implementation
}

// Use structured logging
log := logf.FromContext(ctx)
log.Info("Reconciling RestAPI", "name", restAPI.Name)
```

### YAML Formatting
```yaml
# Use consistent indentation (2 spaces)
apiVersion: apps.aws.com/v1
kind: RestAPI
metadata:
  name: example
spec:
  model:
    enabled: true
    image: "example/model:v1"
    port: 8080
```

### Commit Messages
```bash
# Format: type(scope): description
feat(controller): add blue-green deployment support
fix(api): resolve validation error for health checks
docs(readme): update installation instructions
test(e2e): add auto-scaling test cases
```

## 🔍 Code Review Process

### Before Submitting
- [ ] Tests pass locally
- [ ] Code follows style guidelines
- [ ] Documentation updated
- [ ] Examples work correctly
- [ ] No breaking changes (or documented)

### PR Requirements
- [ ] Clear description of changes
- [ ] Link to related issues
- [ ] Test coverage maintained/improved
- [ ] Documentation updated
- [ ] Examples provided for new features

### Review Checklist
- [ ] **Functionality**: Does it work as intended?
- [ ] **Performance**: No performance regressions
- [ ] **Security**: No security vulnerabilities
- [ ] **Maintainability**: Code is readable and maintainable
- [ ] **Compatibility**: Backward compatibility preserved

## 🐛 Bug Reports

### Issue Template
```markdown
**Bug Description**
Clear description of the bug

**Steps to Reproduce**
1. Deploy operator
2. Apply RestAPI resource
3. Observe error

**Expected Behavior**
What should happen

**Actual Behavior**
What actually happens

**Environment**
- Kubernetes version:
- Operator version:
- Go version:

**Logs**
```
kubectl logs deployment/amazonq-golang-restfulapi-operator-controller-manager
```

**Additional Context**
Any other relevant information
```

## 💡 Feature Requests

### Enhancement Template
```markdown
**Feature Description**
Clear description of the proposed feature

**Use Case**
Why is this feature needed?

**Proposed Solution**
How should this be implemented?

**Alternatives Considered**
Other approaches considered

**Additional Context**
Any other relevant information
```

## 🏗️ Architecture Guidelines

### Adding New Features
1. **Design First**: Create design document
2. **API Changes**: Update CRD if needed
3. **Controller Logic**: Implement reconciliation
4. **Tests**: Add comprehensive tests
5. **Documentation**: Update all relevant docs

### Code Organization
```
api/v1/                 # CRD definitions
internal/controller/    # Controller logic
config/                 # Kubernetes manifests
examples/              # Usage examples
docs/                  # Documentation
test/                  # Test files
```

### Dependencies
- Minimize external dependencies
- Use standard Kubernetes libraries
- Pin dependency versions
- Document dependency changes

## 🚀 Release Process

### Version Management
- Follow semantic versioning (v1.2.3)
- Tag releases in git
- Update CHANGELOG.md
- Build and publish container images

### Release Checklist
- [ ] All tests pass
- [ ] Documentation updated
- [ ] Examples verified
- [ ] Breaking changes documented
- [ ] Migration guide provided (if needed)

## 📞 Getting Help

### Communication Channels
- **Issues**: GitHub issues for bugs/features
- **Discussions**: GitHub discussions for questions
- **Slack**: #restapi-operator channel
- **Email**: maintainers@example.com

### Maintainers
- **@maintainer1**: Core features, API design
- **@maintainer2**: Testing, CI/CD
- **@maintainer3**: Documentation, examples

## 🎯 Good First Issues

### Beginner-Friendly Tasks
- [ ] Add more examples in `examples/` directory
- [ ] Improve error messages and logging
- [ ] Add unit tests for utility functions
- [ ] Update documentation with screenshots
- [ ] Fix typos and improve code comments

### Intermediate Tasks
- [ ] Implement custom metrics for auto-scaling
- [ ] Add support for init containers
- [ ] Enhance blue-green deployment strategies
- [ ] Add webhook validation
- [ ] Implement operator metrics

### Advanced Tasks
- [ ] Multi-cluster support
- [ ] Custom resource validation
- [ ] Advanced scheduling policies
- [ ] Integration with service mesh
- [ ] Performance optimizations

## 📜 License

By contributing, you agree that your contributions will be licensed under the Apache License 2.0.

---

**Thank you for contributing to RestAPI Operator! 🎉**