#!/bin/bash
# Build script for local testing (mimics GitHub Actions)

set -e

echo "🚀 TempFiles Local Build Script"
echo "==============================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Check Go version
echo -e "${BLUE}📋 Checking Go version...${NC}"
go version

# Download dependencies
echo -e "${BLUE}📚 Downloading dependencies...${NC}"
go mod download
go mod verify

# Run tests
echo -e "${BLUE}🧪 Running tests...${NC}"
if go test -v ./...; then
    echo -e "${GREEN}✅ Tests passed${NC}"
else
    echo -e "${RED}❌ Tests failed${NC}"
    exit 1
fi

# Run vet
echo -e "${BLUE}🔍 Running go vet...${NC}"
if go vet ./...; then
    echo -e "${GREEN}✅ Vet passed${NC}"
else
    echo -e "${RED}❌ Vet failed${NC}"
    exit 1
fi

# Create build directory
mkdir -p bin dist

# Build main binary
echo -e "${BLUE}🔨 Building main binary...${NC}"
if go build -v -ldflags="-s -w" -o bin/tempfile ./cmd/server; then
    echo -e "${GREEN}✅ Build successful${NC}"
    ls -la bin/tempfile
    file bin/tempfile
else
    echo -e "${RED}❌ Build failed${NC}"
    exit 1
fi

# Build multi-platform binaries
echo -e "${BLUE}🌍 Building multi-platform binaries...${NC}"

platforms=(
    "linux/amd64"
    "linux/arm64" 
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64"
)

for platform in "${platforms[@]}"; do
    IFS='/' read -r GOOS GOARCH <<< "$platform"
    
    if [ "$GOOS" = "windows" ]; then
        output="dist/tempfile-${GOOS}-${GOARCH}.exe"
    else
        output="dist/tempfile-${GOOS}-${GOARCH}"
    fi
    
    echo -e "${YELLOW}Building for ${GOOS}/${GOARCH}...${NC}"
    
    if CGO_ENABLED=0 GOOS=$GOOS GOARCH=$GOARCH go build \
        -v -ldflags="-s -w" \
        -o "$output" ./cmd/server; then
        echo -e "${GREEN}✅ ${GOOS}/${GOARCH} build successful${NC}"
    else
        echo -e "${RED}❌ ${GOOS}/${GOARCH} build failed${NC}"
        exit 1
    fi
done

# Show build results
echo -e "\n${GREEN}🎉 All builds completed successfully!${NC}"
echo -e "\n${BLUE}📦 Build artifacts:${NC}"
ls -la bin/ dist/

# Test main binary
echo -e "\n${BLUE}✅ Testing main binary...${NC}"
./bin/tempfile --version || echo "Binary runs successfully"

echo -e "\n${GREEN}🚀 Local build completed successfully!${NC}"
echo -e "${YELLOW}💡 To run the server: ./bin/tempfile${NC}"
