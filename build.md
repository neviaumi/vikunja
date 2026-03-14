# Cross-Platform Deployment Guide

This document outlines the deployment process for Vikunja, specifically targeting resource-constrained environments like a Raspberry Pi 3B (1GB RAM). It uses cross-compilation to bypass hardware limitations and employs explicit Git commit hashes for version control alongside the `latest` tag.

## Prerequisites

1.  **Host Machine**: A modern, relatively powerful computer (e.g., PC or Mac) with Docker and the Docker Buildx plugin installed.
2.  **Target Device**: A Raspberry Pi (or similar device) running Docker and Docker Compose.
3.  **Network**: SSH access between the host machine and the target device.

## Step 1: Cross-Compilation

To avoid Out of Memory (OOM) errors on the target device, build the Docker image on your host machine. We tag the image with both the `latest` tag and a short Git commit hash for explicit version tracking.

1.  Navigate to the root of the Vikunja repository on your host machine.
2.  Retrieve the short Git commit hash to use as an explicit version tag:
    ```bash
    export COMMIT_HASH=$(git rev-parse --short HEAD)
    echo "Building version: $COMMIT_HASH"
    ```
3.  Build the image using Docker Buildx, targeting the `linux/arm64` architecture, and tag it with both `latest` and the commit hash:
    ```bash
    docker buildx build --platform linux/arm64 \
      -t vikunja-pi:latest \
      -t vikunja-pi:$COMMIT_HASH \
      --load .
    ```

## Step 2: Archiving the Image

Export the built image into a tar archive so it can be transferred securely over the local network without needing a public container registry.

```bash
docker save vikunja-pi:latest vikunja-pi:$COMMIT_HASH -o vikunja-pi-$COMMIT_HASH.tar
```
*Note: Saving both tags ensures the target device recognizes both the explicit version and the latest pointer when imported.*

## Step 3: Transfer & Load

Transfer the archive to the target device and load it into its local Docker daemon.

1.  Transfer the file via `scp` (replace `<PI_IP_ADDRESS>` and `<USER>` with your actual details):
    ```bash
    scp vikunja-pi-$COMMIT_HASH.tar <USER>@<PI_IP_ADDRESS>:/tmp/
    ```
2.  SSH into your target device:
    ```bash
    ssh <USER>@<PI_IP_ADDRESS>
    ```
3.  Load the imported archive into the Docker daemon:
    ```bash
    docker load -i /tmp/vikunja-pi-$COMMIT_HASH.tar
    ```

## Step 4: Configuration & Deployment

Update your `docker-compose.yml` to specify the explicit commit hash tag. This guarantees the exact version you expect is deployed.

1.  Edit `docker-compose.yml` on the target device:
    ```yaml
    services:
      vikunja:
        image: vikunja-pi:<COMMIT_HASH> # Replace <COMMIT_HASH> with your actual short hash
        environment:
          VIKUNJA_SERVICE_PUBLICURL: http://todo.home.pi/
          VIKUNJA_SERVICE_JWTSECRET: de11dbeff989d29907a547bd57038c6eeca4a197c41e12d705d504c383679f59
          VIKUNJA_DATABASE_PATH: /db/vikunja.db
        ports:
          - 3456:3456
        volumes:
          - ./files:/app/vikunja/files
          - ./config.yml:/app/vikunja/config.yml
          - ./db:/db
        restart: unless-stopped
    ```
2.  Apply the changes and start the container:
    ```bash
    docker compose up -d
    ```

## Rollback Procedures

By tagging images with explicit Git commit hashes alongside `latest`, you retain a history of specific builds directly on your target device.

If a new deployment acts unexpectedly, you can instantly revert to the previous working state:
1.  Open `docker-compose.yml`.
2.  Change the `image: vikunja-pi:<NEW_HASH>` back to `image: vikunja-pi:<PREVIOUS_HASH>`.
3.  Run `docker compose up -d`.

Because the old image is securely tagged with its commit hash, Docker will immediately switch to the previous container without requiring a rebuild or re-transfer.
