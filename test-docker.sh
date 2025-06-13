#!/bin/bash

# Docker Health Check Test Script
echo "🐳 Testing TempFiles Docker Setup"
echo "================================="

# Test 1: Build image
echo "1️⃣ Building Docker image..."
if docker build -t tempfile:test .; then
    echo "✅ Docker build successful"
else
    echo "❌ Docker build failed"
    exit 1
fi

# Test 2: Run container
echo "2️⃣ Starting container..."
CONTAINER_ID=$(docker run -d -p 3001:3000 tempfile:test)
echo "Container ID: $CONTAINER_ID"

# Wait for container to start
echo "⏳ Waiting for container to start..."
sleep 10

# Test 3: Health check
echo "3️⃣ Testing health endpoint..."
if curl -f http://localhost:3001/health; then
    echo "✅ Health check passed"
else
    echo "❌ Health check failed"
fi

# Test 4: Web UI
echo "4️⃣ Testing web UI..."
if curl -f -s http://localhost:3001/ | grep -q "TempFiles"; then
    echo "✅ Web UI accessible"
else
    echo "❌ Web UI not accessible"
fi

# Test 5: Docker health check
echo "5️⃣ Testing Docker internal health check..."
HEALTH_STATUS=$(docker inspect --format='{{.State.Health.Status}}' $CONTAINER_ID)
echo "Health status: $HEALTH_STATUS"

# Cleanup
echo "🧹 Cleaning up..."
docker stop $CONTAINER_ID
docker rm $CONTAINER_ID
docker rmi tempfile:test

echo "================================="
echo "✅ Docker tests completed!"
