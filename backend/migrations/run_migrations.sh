#!/bin/bash
# Simple migration runner for Telegraph

# Database connection from .env
source ../.env

echo "Running migrations for Telegraph..."

# Run each migration in order
for migration in $(ls -1 *.sql | sort); do
    echo "Applying $migration..."
    psql "$DATABASE_URL" < "$migration"
    if [ $? -eq 0 ]; then
        echo "✓ $migration applied successfully"
    else
        echo "✗ $migration failed"
        exit 1
    fi
done

echo "✓ All migrations applied successfully!"
