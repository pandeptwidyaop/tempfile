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

### Option 2: Docker Compose with Nginx

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
    volumes:
      - ./uploads:/app/uploads
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

| Variable | Description | Default | Production Example |
|----------|-------------|---------|-------------------|
| `PUBLIC_URL` | Public URL for download links | `http://localhost:3000` | `https://files.yourdomain.com` |
| `PORT` | Server port | `3000` | `3000` |
| `APP_ENV` | Environment mode | `production` | `production` |
| `DEBUG` | Debug logging | `false` | `false` |
| `CORS_ORIGINS` | Allowed CORS origins | `*` | `*` (public API) or `https://files.yourdomain.com` (web only) |

## Important Notes

1. **PUBLIC_URL is critical** - All download links will use this URL
2. **CORS_ORIGINS** - Set to your actual domain in production
3. **File uploads** - Ensure nginx `client_max_body_size` matches app config
4. **SSL/HTTPS** - Always use HTTPS in production
5. **Headers** - Nginx should pass `X-Forwarded-Proto` for HTTPS detection

## Testing Production Config

```bash
# Test health endpoint
curl https://files.yourdomain.com/health

# Test upload
curl -X POST -F "file=@test.txt" https://files.yourdomain.com/

# Verify download URL in response uses correct domain
```
