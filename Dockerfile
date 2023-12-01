FROM node:alpine

WORKDIR /usr/src/app

COPY package*.json ./

COPY . .

#COPY entrypoint.sh /
#
#RUN chmod +x /entrypoint.sh

ENV PORT 5000

EXPOSE 5000

RUN ls

RUN npm install -g npm@8.3.0

RUN npm install

RUN npm run build

CMD ["npm", "run", "start"]
