#!/bin/bash

# Pull the pre-built image
docker pull ghcr.io/zeeno-atl/claude-code:latest

# Run with current directory
docker run -it --rm -v "$(pwd):/app" ghcr.io/zeeno-atl/claude-code:latest

# Run with API key
docker run -it --rm -v "$(pwd):/app" -e ANTHROPIC_API_KEY="your_api_key" ghcr.io/zeeno-atl/claude-code:latest
