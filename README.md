## pulseUp

Seamless log monitoring for Docker containers with intelligent
action logs for next-level performance and insight.

## Pending Feature
- Insight Log
- Action Log
- Authentication
  
## Getting Started
Pull the latest release with:

```
docker pull mehranzand/pulseup:latest
```

### Using pulseUp with Docker Container:

The easiest way to utilize pulseUp is by running it within a Docker container. Simply follow these steps:
Run the Docker container with the following command:

```
docker run --name pulseup -d --volume=/var/run/docker.sock:/var/run/docker.sock:ro -p 7070:7070 mehranzand/pulseup:latest
```

This command creates a container named "pulseup" from the latest pulseUp image, mounting the Docker Unix socket with read-only permissions (--volume=/var/run/docker.sock:/var/run/docker.sock:ro) and exposing pulseUp on port 7070.

### Using pulseUp with Docker Compose:

Alternatively, you can use Docker Compose to manage your pulseUp service. Here's a sample Docker Compose file:

```
version: "3.4"
services:
  pulseup:
    container_name: pulseup
    image: mehranzand/pulseup:latest
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock:ro
    ports:
      - 7070:7070
```

This configuration achieves the same setup as the Docker command above, allowing you to manage your pulseUp instance in a more structured manner with Docker Compose.
