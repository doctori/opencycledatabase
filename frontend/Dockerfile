FROM node:buster
COPY . /src/
WORKDIR /src
RUN npm install --legacy-peer-deps
RUN npm run build --production
FROM nginx

COPY --from=0 /src/dist/ /usr/share/nginx/html/