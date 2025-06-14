# 🚀 TempFiles Deployment Guide

## Nginx Reverse Proxy Configuration

### Basic Nginx Config (`/etc/nginx/sites-available/tempfiles`)

```nginx
server {
    listen 80;
    server_name files.yourdomain.com;
    
    # Redirect HTTP to HTTPS
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name files.yourdomain.com;
    
    # SSL Configuration
    ssl_certificate /path/to/your/fullchain.pem;
    ssl_certificate_key /path/to/your/privkey.pem;
    
    # Security headers
    add_header X-Frame-Options DENY;
    add_header X-Content-Type-Options nosniff;
    add_header X-XSS-Protection "1; mode=block";
    
    # File upload size limit (should match your app config)
    client_max_body_size 100M;
    
    # Proxy to TempFiles app
    location / {
        proxy_pass http://127.0.0.1:3000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
        
        # Upload timeout
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
        
        # Rate limiting (optional - can also be handled by app)
        # limit_req_zone $binary_remote_addr zone=upload:10m rate=5r/m;
        # limit_req zone=upload burst=10 nodelay;
    }
}
```

## Environment Configuration

### 1. Choose your deployment mode

**Option A: Public API + Web UI**
```bash
cp .env.production.api .env
# API accessible from anywhere + Web UI
```

**Option B: Web UI Only**
```bash
cp .env.production .env
# Only your domain can access API
```

### 2. Update your domain
```bash
# Edit .env and change:
PUBLIC_URL=https://files.yourdomain.com
CORS_ORIGINS=https://files.yourdomain.com,https://yourdomain.com
```

### 3. Build and run
```bash
make build
./bin/tempfile
```

## Docker Deployment

### Option 1: Update existing Dockerfile

Add environment variable support:

```dockerfile
# Add this to your existing Dockerfile
ENV PUBLIC_URL=http://localhost:3000
ENV PORT=3000
ENV APP_ENV=production
ENV DEBUG=false
```

### Option 2: Docker Compose with Nginx and Redis

```yaml
version: '3.8'
services:
  tempfiles:
    build: .
    environment:
      - PUBLIC_URL=https://files.yourdomain.com
      - PORT=3000
      - APP_ENV=production
      - DEBUG=false
      # Rate limiting configuration
      - ENABLE_RATE_LIMIT=true
      - RATE_LIMIT_STORE=redis
      - REDIS_URL=redis://redis:6379
      - RATE_LIMIT_UPLOADS_PER_MINUTE=10
      - RATE_LIMIT_BYTES_PER_HOUR=209715200
      - RATE_LIMIT_TRUSTED_PROXIES=172.16.0.0/12
    volumes:
      - ./uploads:/app/uploads
    depends_on:
      - redis
    restart: unless-stopped
    
  redis:
    image: redis:7-alpine
    command: redis-server --maxmemory 256mb --maxmemory-policy allkeys-lru
    volumes:
      - redis_data:/data
    restart: unless-stopped
    
  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
      - ./ssl:/etc/nginx/ssl
    depends_on:
      - tempfiles
    restart: unless-stopped

volumes:
  redis_data:
```

## Cloud Deployment Examples

### Vercel/Netlify (Static + Serverless)
- Deploy as serverless function
- Set `PUBLIC_URL` in environment variables

### Railway/Render
```bash
# Set environment variables:
PUBLIC_URL=https://your-app.railway.app
PORT=3000
APP_ENV=production
```

### Traditional VPS
```bash
# 1. Clone repository
git clone https://github.com/pandeptwidyaop/tempfile.git
cd tempfile

# 2. Set up environment
cp .env.production .env
# Edit .env with your domain

# 3. Build and run
make build
./bin/tempfile

# 4. Setup systemd service (optional)
sudo systemctl enable tempfiles
sudo systemctl start tempfiles
```

## Environment Variables Reference

### Core Configuration

| Variable | Description | Default | Production Example |
|----------|-------------|---------|-------------------|
| `PUBLIC_URL` | Public URL for download links | `http://localhost:3000` | `https://files.yourdomain.com` |
| `PORT` | Server port | `3000` | `3000` |
| `APP_ENV` | Environment mode | `production` | `production` |
| `DEBUG` | Debug logging | `false` | `false` |
| `CORS_ORIGINS` | Allowed CORS origins | `*` | `*` (public API) or `https://files.yourdomain.com` (web only) |

### Rate Limiting Configuration

| Variable | Description | Default | Production Example |
|----------|-------------|---------|-------------------|
| `ENABLE_RATE_LIMIT` | Enable rate limiting | `false` | `true` |
| `RATE_LIMIT_STORE` | Storage backend | `memory` | `redis` (for distributed) |
| `RATE_LIMIT_UPLOADS_PER_MINUTE` | Max uploads per minute per IP | `5` | `10` |
| `RATE_LIMIT_BYTES_PER_HOUR` | Max bytes per hour per IP | `104857600` | `209715200` (200MB) |
| `RATE_LIMIT_TRUSTED_PROXIES` | Trusted proxy IPs/CIDRs | `127.0.0.1,::1,...` | `172.16.0.0/12,10.0.0.0/8` |
| `RATE_LIMIT_WHITELIST_IPS` | Whitelisted IPs | `` | `192.168.1.0/24,203.0.113.100` |

### Redis Configuration (for distributed rate limiting)

| Variable | Description | Default | Production Example |
|----------|-------------|---------|-------------------|
| `REDIS_URL` | Redis connection URL | `redis://localhost:6379` | `redis://redis-cluster:6379` |
| `REDIS_PASSWORD` | Redis password | `` | `your-secure-password` |
| `REDIS_DB` | Redis database number | `0` | `0` |

## Rate Limiting Deployment Scenarios

### Scenario 1: Single Instance with Memory Store
```bash
# Simple deployment without Redis
ENABLE_RATE_LIMIT=true
RATE_LIMIT_STORE=memory
RATE_LIMIT_UPLOADS_PER_MINUTE=10
RATE_LIMIT_BYTES_PER_HOUR=209715200
```

### Scenario 2: Load Balanced with Redis
```bash
# Multiple instances sharing rate limit data
ENABLE_RATE_LIMIT=true
RATE_LIMIT_STORE=redis
REDIS_URL=redis://redis-cluster:6379
REDIS_PASSWORD=your-password
RATE_LIMIT_TRUSTED_PROXIES=172.16.0.0/12,10.0.0.0/8
```

### Scenario 3: Behind Cloudflare
```bash
# Configure for Cloudflare proxy
ENABLE_RATE_LIMIT=true
RATE_LIMIT_TRUSTED_PROXIES=173.245.48.0/20,103.21.244.0/22,103.22.200.0/22
RATE_LIMIT_IP_HEADERS=CF-Connecting-IP,X-Forwarded-For
RATE_LIMIT_WHITELIST_IPS=your-office-ip/32
```

## Important Notes

1. **PUBLIC_URL is critical** - All download links will use this URL
2. **CORS_ORIGINS** - Set to your actual domain in production
3. **File uploads** - Ensure nginx `client_max_body_size` matches app config
4. **SSL/HTTPS** - Always use HTTPS in production
5. **Headers** - Nginx should pass `X-Forwarded-Proto` and `X-Real-IP` for rate limiting
6. **Rate Limiting** - Use Redis backend for distributed deployments
7. **Trusted Proxies** - Configure properly for accurate IP detection

## Testing Production Config

```bash
# Test health endpoint
curl https://files.yourdomain.com/health

# Test upload
curl -X POST -F "file=@test.txt" https://files.yourdomain.com/

# Verify download URL in response uses correct domain

# Test rate limiting (should return 429 after limit exceeded)
for i in {1..10}; do curl -X POST -F "file=@test.txt" https://files.yourdomain.com/; done

# Check rate limit headers
curl -I -X POST -F "file=@test.txt" https://files.yourdomain.com/
```

## Redis Deployment Notes

### Redis Security
```bash
# Redis configuration for production
redis-server --requirepass your-secure-password \
             --maxmemory 512mb \
             --maxmemory-policy allkeys-lru \
             --save 900 1 \
             --appendonly yes
```

### Redis Monitoring
```bash
# Monitor Redis for rate limiting
redis-cli monitor | grep ratelimit

# Check rate limit keys
redis-cli keys "ratelimit:*"

# Monitor memory usage
redis-cli info memory
```
