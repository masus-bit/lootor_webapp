FROM node:alpine

WORKDIR /usr/src/app

COPY package*.json ./

COPY . .

ENV PORT 5000

RUN npm install -g npm@8.3.0

RUN npm install

RUN npm run build

RUN pwd

RUN ls -la

RUN cat package.json

CMD ["npm", "run", "start"]
