# 🐳 Docker Deployment Guide

This guide covers deploying Tempfile using Docker images from GitHub Container Registry (GHCR).

## Quick Start

```bash
# Pull and run the latest stable version
docker pull ghcr.io/your-username/tempfile:latest
docker run -d --name tempfile -p 3000:3000 -v $(pwd)/uploads:/app/uploads ghcr.io/your-username/tempfile:latest
```

## Available Images

### Production Images
- `ghcr.io/your-username/tempfile:latest` - Latest stable release
- `ghcr.io/your-username/tempfile:vX.X.X` - Specific versions

### Development Images  
- `ghcr.io/your-username/tempfile:develop` - Latest development build
- `ghcr.io/your-username/tempfile:feature-*` - Feature branch builds

## Docker Compose

```yaml
version: '3.8'

services:
  tempfile:
    image: ghcr.io/your-username/tempfile:latest
    container_name: tempfile
    restart: unless-stopped
    ports:
      - "3000:3000"
    volumes:
      - ./uploads:/app/uploads
    environment:
      - PORT=3000
      - APP_ENV=production
      - PUBLIC_URL=https://files.yourdomain.com
    healthcheck:
      test: ["CMD", "wget", "--spider", "http://localhost:3000/health"]
      interval: 30s
      timeout: 10s
      retries: 3
```

Run with: `docker-compose up -d`

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `3000` | Application port |
| `APP_ENV` | `production` | Environment mode |
| `DEBUG` | `false` | Debug logging |
| `PUBLIC_URL` | `http://localhost:3000` | Public URL |
| `ENABLE_WEB_UI` | `true` | Enable web interface |

## Production Deployment

### With Reverse Proxy

```yaml
version: '3.8'

services:
  tempfile:
    image: ghcr.io/your-username/tempfile:latest
    restart: unless-stopped
    volumes:
      - uploads_data:/app/uploads
    environment:
      - PUBLIC_URL=https://files.yourdomain.com

  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf:ro
      - ./ssl:/etc/nginx/ssl:ro
    depends_on:
      - tempfile

volumes:
  uploads_data:
```

### Security Hardening

```bash
docker run -d \
  --name tempfile \
  --security-opt no-new-privileges:true \
  --cap-drop ALL \
  --read-only \
  --tmpfs /tmp \
  --memory="512m" \
  --cpus="1.0" \
  -v uploads:/app/uploads \
  -p 3000:3000 \
  ghcr.io/your-username/tempfile:latest
```

## Monitoring

```bash
# Check health
curl http://localhost:3000/health

# View logs
docker logs -f tempfile

# Monitor resources
docker stats tempfile
```

## Updates

```bash
# Pull latest version
docker pull ghcr.io/your-username/tempfile:latest

# Recreate container
docker-compose up -d --force-recreate
```

For more detailed information, see the full [Docker GHCR Documentation](./docs/DOCKER_GHCR.md).
