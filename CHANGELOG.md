# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.1] - 2025-09-20

### Added
- Initial RestAPI Operator implementation
- MVC+R pattern support (Model, View, Controller, Repository)
- Auto-scaling with HorizontalPodAutoscaler
- Health monitoring with configurable probes
- Blue-green deployment strategy
- Comprehensive documentation and examples

### Changed
- N/A

### Deprecated
- N/A

### Removed
- N/A

### Fixed
- N/A

### Security
- RBAC configuration for operator permissions

## [1.0.0] - 2025-09-20

### Added
- RestAPI Custom Resource Definition (CRD)
- RestAPI Controller with reconciliation logic
- Blue-green deployment manager
- Support for individual MVC+R component configuration
- Auto-scaling based on CPU and memory metrics
- Health check configuration for all components
- Service creation for component discovery
- Environment variable management
- Deployment scripts and examples
- Complete documentation suite:
  - HOW-TO-USE.md - Comprehensive usage guide
  - TUTORIAL.md - Hands-on learning tutorial
  - QUICK-START.md - 5-minute setup guide
  - ARCHITECTURE.md - System design documentation
  - CONTRIBUTOR.md - Developer contribution guide

### Technical Details
- Built with operator-sdk and kubebuilder
- Go 1.21+ support
- Kubernetes 1.20+ compatibility
- Apache 2.0 license