FROM node:24-alpine

WORKDIR /app

RUN addgroup -g 1001 -S nodejs && \
    adduser -S nodejs -u 1001

COPY package.json pnpm-workspace.yaml ./
COPY services/frontend/package.json ./services/frontend/
COPY pnpm-lock.yaml* ./

RUN corepack prepare pnpm@8.15.0 --activate
RUN corepack enable pnpm && \
    pnpm install --frozen-lockfile --filter=@jira-clone/frontend

COPY services/frontend/ ./services/frontend/

RUN chown -R nodejs:nodejs /app
USER nodejs

EXPOSE 5173

CMD ["pnpm", "run", "dev"]