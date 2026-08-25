# Evidencias

## 1. Protección de main

Se intentó hacer un push directamente a main. GitHub rechazó el cambio porque la rama está protegida y los cambios deben realizarse mediante un Pull Request.

![Push rechazado](img/01-push-rechazado.png)

## 2. Pull Request

Se creó un Pull Request desde feature/seccion-instalacion hacia main para incorporar los cambios realizados.

![Pull Request](img/02-conflicto-PR.png)

## 3. Resolución de conflicto

Se generó un conflicto porque dos ramas modificaron el mismo contenido de manera diferente. Se resolvió manualmente eliminando los marcadores de conflicto 

![Conflicto](img/03-marcadores-conflicto.png)

## 4. Release

Se publicó la versión v1.0.0 como Release final del TP1.

![Release publicada](img/04-release-publicada.png)


## TP2 — Contenedores

### 5. Sistema completo funcionando (`docker compose up`)

Los tres servicios (base de datos, backend y frontend) levantados con `docker compose up -d --build`, con la base de datos en estado `healthy`.

[![Sistema funcionando](https://github.com/canderojo/ingsoft3-tp01/raw/main/img/05-compose-up.png)](/canderojo/ingsoft3-tp01/blob/main/img/05-compose-up.png)

### 6. Prueba de persistencia

Se creó un turno y se verificó que sobrevive a un `docker compose down` / `up` (el volumen conserva los datos), y que `docker compose down -v` sí lo borra (elimina también el volumen).

[![Prueba de persistencia](https://github.com/canderojo/ingsoft3-tp01/raw/main/img/06-persistencia.png)](/canderojo/ingsoft3-tp01/blob/main/img/06-persistencia.png)

### 7. Comparación de tamaño: imagen final vs. imagen de build

Comparación entre el tamaño de la imagen de compilación (`golang:1.25-alpine`, con el compilador incluido) y la imagen final del backend (`turnos-backend:dev`, solo el binario compilado). El multi-stage build reduce el tamaño final a una fracción del de la imagen de compilación.

[![Comparación de tamaños](https://github.com/canderojo/ingsoft3-tp01/raw/main/img/07-comparacion-tamanos.png)](/canderojo/ingsoft3-tp01/blob/main/img/07-comparacion-tamanos.png)

### 8. Imágenes publicadas en el registry

Las imágenes `turnos-backend` y `turnos-frontend` publicadas en GitHub Container Registry (ghcr.io), con visibilidad pública, verificado con un `docker pull` sin estar autenticada.

[![Imágenes publicadas](https://github.com/canderojo/ingsoft3-tp01/raw/main/img/08-registry-publicado.png)](/canderojo/ingsoft3-tp01/blob/main/img/08-registry-publicado.png)