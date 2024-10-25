FROM node:21-alpine AS base

WORKDIR /app

COPY package.json package-lock.json yarn.lock ./

RUN npm install

COPY . .

RUN npm run build

EXPOSE 3000

CMD npm run start
