# Use Node.js base image for building the React app
FROM node:20-alpine AS builder

# Set working directory
WORKDIR /app

# Copy package files
COPY package*.json ./

# Install ALL dependencies (including devDependencies for build tools)
RUN npm ci

# Copy source code
COPY . .

# Build the application
RUN npm run build

# Production stage - serve the built app using a lightweight web server
FROM node:20-alpine

# Install a simple web server to serve static files
RUN npm install -g serve

# Create a non-root user
RUN addgroup -S ecom && \
    adduser -S ecom -G ecom -u 1000

# Set working directory and ownership
WORKDIR /app
RUN chown ecom:ecom /app

# Use the non-root user
USER ecom

# Copy the built application from the builder stage
COPY --from=builder /app/dist ./dist

# Expose port 3000 (or whatever port your app uses)
EXPOSE 3000

# Serve the app
CMD ["serve", "-s", "dist", "-l", "3000"]
