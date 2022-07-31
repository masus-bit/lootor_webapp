FROM node:slim

WORKDIR /usr/src/app

COPY package*.json ./

ENV PORT 5000

RUN npm install -g npm@8.3.0

RUN npm install

COPY . .

RUN npm run build

CMD ["npm", "run", "start"]
