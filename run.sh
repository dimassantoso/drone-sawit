#!/usr/bin/env bash
make init                      # Initializes 
make                           # Builds the binary
make test                      # Runs unit tests with coverage 
docker compose up --build -d   # Runs the docker to start the API and database.
sleep 30                       # Wait until API runs
make test_api                  # Runs the API testing with data.
docker compose down --volumes  # Stops the docker containers and removes the volumes.