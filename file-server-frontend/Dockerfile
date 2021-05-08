FROM node:latest
WORKDIR /app
ENV PATH /app/node_modules/.bin:$PATH
COPY package.json /app/package.json
RUN npm install -g serve
RUN npm install
RUN npm install -g @angular/cli@7.3.9
COPY . /app
RUN ng build
CMD serve -s dist/file-server-frontend
