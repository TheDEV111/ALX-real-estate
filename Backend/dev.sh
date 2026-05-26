#!/bin/bash

# Backend development script

set -e

echo "🚀 Starting ALX Real Estate Backend..."

# Check if docker-compose is available
if ! command -v docker-compose &> /dev/null; then
    echo "❌ docker-compose not found. Please install Docker and docker-compose."
    exit 1
fi

# Start services
docker-compose up --build -d

echo "✅ Backend running on http://localhost:8080"
echo "✅ PostgreSQL running on localhost:5432"
echo ""
echo "📋 Useful commands:"
echo "  docker-compose logs -f backend    # View backend logs"
echo "  docker-compose logs -f postgres   # View database logs"
echo "  docker-compose down               # Stop services"
echo "  docker-compose exec backend sh    # Access backend container"