FROM node:20-alpine AS builder

ARG NEXT_PUBLIC_API_URL=http://localhost:8080

WORKDIR /app
COPY package*.json ./
RUN npm install

COPY . .
ENV NEXT_PUBLIC_API_URL=${NEXT_PUBLIC_API_URL}
RUN npm run build

# Production image
FROM node:20-alpine

WORKDIR /app
COPY package*.json ./
RUN npm install --production

COPY --from=builder /app/.next ./.next
COPY --from=builder /app/public ./public
COPY --from=builder /app/next.config.ts ./next.config.ts
COPY --from=builder /app/package.json ./package.json

EXPOSE 3000
CMD ["npm", "start"]
