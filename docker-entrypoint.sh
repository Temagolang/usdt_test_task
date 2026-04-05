#!/bin/sh
set -e

MAX_RETRIES=30

auto_migrate() {
  echo "entrypoint: applying migrations..."
  attempt=0
  while [ "$attempt" -lt "$MAX_RETRIES" ]; do
    if ./app migrate up 2>&1; then
      echo "entrypoint: migrations applied"
      return 0
    fi
    attempt=$((attempt + 1))
    echo "entrypoint: migrate up failed, retry $attempt/$MAX_RETRIES..."
    sleep 1
  done
  echo "entrypoint: failed to apply migrations after $MAX_RETRIES attempts" >&2
  exit 1
}

# Determine whether to auto-migrate based on command.
# ./app or ./app grpc  → migrate + start
# ./app migrate ...    → pass through (user manages migrations)
# anything else        → pass through
case "$*" in
  "./app"|"./app grpc")
    auto_migrate
    exec "$@"
    ;;
  "./app migrate"*)
    exec "$@"
    ;;
  *)
    exec "$@"
    ;;
esac