FROM node:alpine

WORKDIR /usr/src/app

COPY package*.json ./

COPY . .

#COPY entrypoint.sh /
#
#RUN chmod +x /entrypoint.sh

ENV PORT 5000

EXPOSE 5000

RUN cat .env

RUN npm install -g npm@10.2.4

RUN npm install

RUN npm run build

CMD ["npm", "run", "start"]
