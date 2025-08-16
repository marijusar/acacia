FROM node:24-alpine

WORKDIR /app

RUN addgroup -g 1001 -S nodejs && \
    adduser -S nodejs -u 1001

COPY package.json pnpm-workspace.yaml ./
COPY services/dashboard-api/package.json ./services/dashboard-api/
COPY pnpm-lock.yaml* ./

RUN corepack prepare pnpm@8.15.0 --activate
RUN corepack enable pnpm && \
    pnpm install --frozen-lockfile --filter=dashboard-api

# COPY services/dashboard-api/ ./ # Commented out for development - files are mounted via volume

RUN chown -R nodejs:nodejs /app
USER nodejs

EXPOSE 3000

CMD ["pnpm", "run", "dev"]
