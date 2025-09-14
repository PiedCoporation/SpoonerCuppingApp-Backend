#!/bin/bash

# Coffee Cupping Backend Deployment Script
# Usage: ./deploy.sh [dev|prod]

set -e

ENVIRONMENT=${1:-dev}

if [[ "$ENVIRONMENT" != "dev" && "$ENVIRONMENT" != "prod" ]]; then
    echo "❌ Error: Environment must be 'dev' or 'prod'"
    echo "Usage: ./deploy.sh [dev|prod]"
    exit 1
fi

echo "🚀 Deploying to $ENVIRONMENT environment..."

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Error: Docker is not running"
    exit 1
fi

# Set project name based on environment
PROJECT_NAME="coffee-cupping-$ENVIRONMENT"

# Build and deploy based on environment
if [[ "$ENVIRONMENT" == "dev" ]]; then
    echo "📦 Building and deploying development environment..."
    
    # Stop existing containers
    echo "Stopping existing containers..."
    docker compose -f docker-compose.yml -f docker-compose.dev.yml -p $PROJECT_NAME down || echo "No existing containers to stop"
    
    # Build and start services
    echo "Building and starting services..."
    docker compose -f docker-compose.yml -f docker-compose.dev.yml -p $PROJECT_NAME up -d --build
    
    echo "⏳ Waiting for services to be ready..."
    sleep 10
    
    echo "🔍 Checking service health..."
    if curl -f http://localhost:8080/health > /dev/null 2>&1; then
        echo "✅ Development deployment successful!"
        echo "🌐 API is available at: http://localhost:8080"
    else
        echo "❌ Health check failed"
        docker compose -f docker-compose.yml -f docker-compose.dev.yml -p $PROJECT_NAME logs
        exit 1
    fi
    
elif [[ "$ENVIRONMENT" == "prod" ]]; then
    echo "📦 Building and deploying production environment..."
    
    # Check if required environment variables are set
    required_vars=("POSTGRES_PASSWORD" "JWT_ACCESS_TOKEN_KEY" "JWT_REFRESH_TOKEN_KEY" "JWT_REGISTER_TOKEN_KEY" "JWT_LOGIN_TOKEN_KEY" "SMTP_HOST" "SMTP_USERNAME" "SMTP_APP_PASSWORD")
    
    for var in "${required_vars[@]}"; do
        if [[ -z "${!var}" ]]; then
            echo "❌ Error: Environment variable $var is not set"
            exit 1
        fi
    done
    
    # Stop existing containers
    echo "Stopping existing containers..."
    docker compose -f docker-compose.yml -f docker-compose.prod.yml -p $PROJECT_NAME down || echo "No existing containers to stop"
    
    # Build and start services
    echo "Building and starting services..."
    docker compose -f docker-compose.yml -f docker-compose.prod.yml -p $PROJECT_NAME up -d --build
    
    echo "⏳ Waiting for services to be ready..."
    sleep 15
    
    echo "🔍 Checking service health..."
    if curl -f http://localhost:8080/health > /dev/null 2>&1; then
        echo "✅ Production deployment successful!"
        echo "🌐 API is available at: http://localhost:8080"
    else
        echo "❌ Health check failed"
        docker compose -f docker-compose.yml -f docker-compose.prod.yml -p $PROJECT_NAME logs
        exit 1
    fi
fi

echo "📊 Deployment Summary:"
echo "Environment: $ENVIRONMENT"
echo "Project Name: $PROJECT_NAME"
echo "Timestamp: $(date)"
echo "Docker Compose Files: docker-compose.yml, docker-compose.$ENVIRONMENT.yml"

# Show running containers
echo "🐳 Running containers:"
docker compose -f docker-compose.yml -f docker-compose.$ENVIRONMENT.yml -p $PROJECT_NAME ps
