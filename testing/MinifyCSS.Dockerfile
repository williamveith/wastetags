# node:23.5-alpine3.21
FROM node@sha256:c61b6b12a3c96373673cd52d7ecee2314e82bca5d541eecf0bc6aee870c8c6f7

RUN npm install -g clean-css-cli@5.6.3

WORKDIR /app

COPY ./build/tools/minifycss.sh /app/minifycss.sh

RUN chmod +x /app/minifycss.sh

ENTRYPOINT ["/app/minifycss.sh"]