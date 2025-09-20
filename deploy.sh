#!/bin/bash

set -e

echo "🚀 Deploying RestAPI Operator for MVC+R Pattern Management"

# Check if kubectl is available
if ! command -v kubectl &> /dev/null; then
    echo "❌ kubectl is not installed or not in PATH"
    exit 1
fi

# Check if cluster is accessible
if ! kubectl cluster-info &> /dev/null; then
    echo "❌ Cannot connect to Kubernetes cluster"
    exit 1
fi

echo "✅ Kubernetes cluster is accessible"

# Install CRDs
echo "📦 Installing Custom Resource Definitions..."
make install

# Deploy the operator
echo "🔧 Deploying RestAPI Operator..."
make deploy

# Wait for operator to be ready
echo "⏳ Waiting for operator to be ready..."
kubectl wait --for=condition=available --timeout=300s deployment/amazonq-golang-restfulapi-operator-controller-manager -n amazonq-golang-restfulapi-operator-system

echo "✅ RestAPI Operator deployed successfully!"

# Deploy sample RestAPI
echo "📋 Deploying sample RestAPI..."
kubectl apply -f config/samples/apps_v1_restapi.yaml

echo "🎉 Deployment complete!"
echo ""
echo "To check the status:"
echo "  kubectl get restapi -A"
echo "  kubectl get pods -n default"
echo "  kubectl get services -n default"
echo ""
echo "To view logs:"
echo "  kubectl logs -f deployment/amazonq-golang-restfulapi-operator-controller-manager -n amazonq-golang-restfulapi-operator-system"