FROM node:21-alpine AS base

# Build argument for store service URL (optional)
ARG STORE_SERVICE_URL

WORKDIR /app

COPY package.json package-lock.json yarn.lock ./

RUN npm install

COPY . .

RUN if [ -n "$STORE_SERVICE_URL" ]; then \
      echo "✓ Build argument provided, replacing .env file"; \
      echo "STORE_SERVICE_URL=${STORE_SERVICE_URL}" > .env; \
    else \
      echo "✓ No build argument, using existing .env file"; \
    fi

RUN npm run build

EXPOSE 3000

CMD npm run start
