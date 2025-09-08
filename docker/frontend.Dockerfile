FROM node:24-alpine

WORKDIR /app

RUN addgroup -g 1001 -S nodejs && \
    adduser -S nodejs -u 1001

COPY services/frontend/package.json ./
COPY services/frontend/package-lock.json* ./
COPY services/frontend/pnpm-lock.yaml* ./

RUN corepack prepare pnpm@10.5.2 --activate
RUN corepack enable pnpm && \
    pnpm install --frozen-lockfile

COPY services/frontend/ ./

RUN chown -R nodejs:nodejs /app
USER nodejs

EXPOSE 5173

CMD ["pnpm", "run", "dev"]
