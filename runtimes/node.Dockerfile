FROM node:12-alpine

RUN mkdir -p /usr/app
WORKDIR /usr/app

COPY package*.json ./
RUN npm install

COPY . .

EXPOSE 8000
CMD [ "node", "index.js" ]
