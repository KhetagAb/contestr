# Use an official Node.js runtime as a parent image
FROM oven/bun:1 AS build

WORKDIR /app/contestr-front

COPY contestr-front/package*.json ./
RUN bun install

COPY contestr-front/ .
COPY contestr/api/openapi.yaml ../contestr/api/openapi.yaml

RUN bun run build

# Use a lightweight web server to serve the built files
FROM nginx:alpine

COPY --from=build /app/contestr-front/dist /opt/build
COPY contestr-front/nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
