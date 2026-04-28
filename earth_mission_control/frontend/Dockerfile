FROM node:20-alpine AS build
WORKDIR /app

ARG VITE_EARTH_BACKEND_URL=http://localhost:8083
ENV VITE_EARTH_BACKEND_URL=${VITE_EARTH_BACKEND_URL}

COPY package.json tsconfig.json vite.config.ts ./
COPY index.html ./
COPY src ./src

RUN npm install
RUN npm run build

FROM nginx:1.27-alpine
COPY --from=build /app/dist /usr/share/nginx/html

EXPOSE 80
