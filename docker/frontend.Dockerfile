FROM node:24-alpine

WORKDIR /app

RUN corepack enable pnpm

EXPOSE 5173

CMD ["pnpm", "run", "dev"]
