FROM node:alpine

WORKDIR /usr/src/app

COPY package*.json ./

COPY . .

#COPY entrypoint.sh /
#
#RUN chmod +x /entrypoint.sh

ENV PORT 5111

EXPOSE 5111

RUN cat .env

RUN npm install -g npm@8.3.0

RUN npm install

RUN npm run typeorm migration:run

RUN npm run build

CMD ["npm", "run", "start"]
