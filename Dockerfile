FROM node:8-alpine

COPY main.js index.js
COPY package.json package.json
RUN npm i

CMD [ "node", "index.js" ]
