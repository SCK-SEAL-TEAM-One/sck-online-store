FROM node:21-alpine AS base

# Build argument for store service URL
ARG STORE_SERVICE_URL=http://localhost

WORKDIR /app

COPY package.json package-lock.json yarn.lock ./

RUN npm install

COPY . .

# Pass build argument as environment variable for Next.js build
ENV STORE_SERVICE_URL=${STORE_SERVICE_URL}

RUN npm run build

EXPOSE 3000

CMD npm run start
