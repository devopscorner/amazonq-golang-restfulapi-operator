---
name: Bug report
about: Create a report to help us improve
title: '[BUG] '
labels: bug
assignees: ''
---

## Bug Description
A clear and concise description of what the bug is.

## Steps to Reproduce
1. Deploy operator with `make deploy`
2. Apply RestAPI resource `kubectl apply -f config/samples/apps_v1_restapi.yaml`
3. Observe the error

## Expected Behavior
A clear description of what you expected to happen.

## Actual Behavior
A clear description of what actually happened.

## Environment
- **Kubernetes version**: [e.g. v1.28.0]
- **Operator version**: [e.g. v1.0.0]
- **Go version**: [e.g. 1.21.4]
- **OS**: [e.g. macOS, Linux, Windows]

## Logs
```bash
# Operator logs
kubectl logs deployment/amazonq-golang-restfulapi-operator-controller-manager -n amazonq-golang-restfulapi-operator-system

# RestAPI resource status
kubectl describe restapi <name>

# Pod logs (if applicable)
kubectl logs <pod-name>
```

## RestAPI Resource
```yaml
# Paste your RestAPI resource YAML here
```

## Additional Context
Add any other context about the problem here, such as:
- Screenshots
- Related issues
- Workarounds attempted