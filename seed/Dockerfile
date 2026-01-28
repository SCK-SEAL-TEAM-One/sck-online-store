FROM node:20

WORKDIR /app

COPY package*.json ./

RUN npm install

COPY generate-store-data.js ./

RUN mkdir -p /app/output

CMD ["node", "generate-store-data.js"]