# docker build -t jsmini -f build/JSMini.Dockerfile .
# docker run --rm -it -v $(pwd)/cmd/wastetags/assets/js:/app/input -v $(pwd)/dist/wastetags/assets/js:/app/output jsmini:latest sh
# terser ./input/arrow-navigation.js -o ./output/arrow-navigation.js
# node:23.5-alpine3.21
FROM node@sha256:c61b6b12a3c96373673cd52d7ecee2314e82bca5d541eecf0bc6aee870c8c6f7

RUN npm install -g terser@5.14.2

WORKDIR /app

ENTRYPOINT ["sh", "-c", "for file in ./input/*.js; do filename=$(basename \"$file\"); terser \"$file\" -o ./output/\"$filename\"; done"]