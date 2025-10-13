FROM node:24-alpine

WORKDIR /app

# Copy package files
COPY services/frontend/package.json services/frontend/package-lock.json ./

# Install dependencies in the container with the correct platform binaries
RUN npm ci

# Copy the rest of the application
COPY services/frontend/ ./

EXPOSE 5173

CMD ["npm", "run", "dev"]
