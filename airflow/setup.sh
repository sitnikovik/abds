#!/bin/bash

set -e

echo "=========================================="
echo "🚀 Starting Airflow initialization..."
echo "=========================================="

# Install Python packages
echo "📦 Installing Python packages..."
pip install --no-cache-dir clickhouse-connect requests

# Run database migrations
echo "🗄️  Running database migrations..."
airflow db migrate

# Create admin user (if not exists)
echo "👤 Creating admin user..."
airflow users create \
    -r Admin \
    -u admin \
    -e admin@example.com \
    -f Admin \
    -l User \
    -p admin || echo "⚠️  Admin user already exists"

# Import connections from JSON (if file exists)
CONNECTIONS_FILE="/opt/airflow/config/connections.json"
if [ -f "$CONNECTIONS_FILE" ]; then
    echo "🔗 Importing connections from $CONNECTIONS_FILE..."
    airflow connections import "$CONNECTIONS_FILE" || echo "⚠️  Failed to import connections"
else
    echo "⚠️  Connections file not found: $CONNECTIONS_FILE"
fi

echo "=========================================="
echo "✅ Airflow initialization completed!"
echo "=========================================="
