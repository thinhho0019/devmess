### Multi-stage Dockerfile for React + Vite (production)
### Build with Node, serve static files with nginx (SPA fallback)

# --- Builder stage ---
FROM node:18-alpine AS builder
WORKDIR /app

# Install dependencies (use package-lock if present)
COPY package*.json ./
# copy lock file if exists (npm/yarn/pnpm) - speeds installs in CI
COPY package-lock.json ./package-lock.json

RUN npm ci --silent

# Copy source and build
COPY . .
# ENV NODE_ENV=production
RUN npm run build
CMD ["sh", "-c", "cp -r dist/* /app_dist/ && echo 'Build complete' && tail -f /dev/null"]

# # --- Production stage ---
# FROM nginx:alpine AS production

# # Remove default nginx static content
# RUN rm -rf /usr/share/nginx/html/*

# # Copy built files from builder
# COPY --from=builder /app/dist /usr/share/nginx/html

# # Create a small nginx conf for SPA fallback routing
# RUN rm /etc/nginx/conf.d/default.conf
# RUN mkdir -p /etc/nginx/conf.d
# RUN { \
#   echo 'server {'; \
#   echo '  listen 80;'; \
#   echo '  server_name _;'; \
#   echo '  root /usr/share/nginx/html;'; \
#   echo '  index index.html;'; \
#   echo ''; \
#   echo '  location / {'; \
#   echo '    try_files $uri $uri/ /index.html;'; \
#   echo '  }'; \
#   echo ''; \
#   echo '  # Optional: disable access log to reduce noise'; \
#   echo '  access_log off;'; \
#   echo '}'; \
# } > /etc/nginx/conf.d/default.conf

# # Expose port and add a simple healthcheck
# EXPOSE 80
# HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 CMD wget -qO- http://localhost:80/ >/dev/null || exit 1

# CMD ["/docker-entrypoint.sh", "nginx", "-g", "daemon off;"]
