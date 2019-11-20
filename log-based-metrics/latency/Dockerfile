FROM node:10

# Create app directory
WORKDIR /usr/src/app

# Install app dependencies
COPY package*.json ./
RUN npm install --save sleep
RUN npm install

# Bundle app source
COPY . .

EXPOSE 8080
CMD [ "node", "latency.js" ]