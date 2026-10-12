# ---- Etapa 1: build ----
FROM node:22-alpine AS build
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

# ---- Etapa 2: tests (parte del build: ya tiene Node y las dependencias) ----
FROM build AS test
ENTRYPOINT ["npm", "run", "test:ci"]

# ---- Etapa 3: nginx sirve los estáticos + proxy hacia el backend ----
FROM nginx:alpine
# La plantilla va en templates/ y no en conf.d/: nginx leería conf.d/
# tal cual, con los ${...} sin reemplazar, y no arrancaría. Al iniciar,
# la imagen de nginx reemplaza las variables y genera conf.d/default.conf.
COPY default.conf.template /etc/nginx/templates/default.conf.template
COPY --from=build /app/dist /usr/share/nginx/html
ENV BACKEND_URL=http://backend:8080 \
    DNS_RESOLVER=127.0.0.11
EXPOSE 80