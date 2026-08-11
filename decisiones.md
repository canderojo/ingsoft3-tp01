# Decisiones

## 1. Por qué Git no pudo resolver el conflicto automáticamente

Git no pudo resolver el conflicto automáticamente porque las ramas `feature/titulo-a` y `feature/titulo-b` modificaron el mismo contenido del archivo de maneras diferentes. Git no podía determinar por sí solo cuál de las dos versiones era la correcta.

El conflicto podría haberse evitado si las dos ramas no hubieran modificado las mismas líneas del archivo, o si los cambios se hubieran integrado antes de que las ramas se alejaran demasiado de `main`.

Para resolverlo, revisé las dos versiones, elegí qué contenido debía quedar y eliminé los marcadores `<<<<<<<`, `=======` y `>>>>>>>`.

## 2. Problemas encontrados y cómo los solucioné

El primer problema que surgio fue el intento de hacer `push` directamente a `main`. GitHub lo rechazó porque la rama estaba protegida, por lo que los cambios debían hacerse mediante un Pull Request. Esto permitió comprobar que la protección estaba funcionando correctamente.

Finalmente, al trabajar con dos ramas se produjo un conflicto. Lo resolví desde GitHub revisando las diferencias entre ambas versiones y seleccionando manualmente el contenido que debía conservarse.

## 3. Declaración de uso de IA

Utilicé inteligencia artificial como herramienta de apoyo durante la realización del TP, principalmente para escribir los contenidos de los archivos y avanzar de manera guiada en la configuración del repositorio, ramas, Pull Requests, resolución de conflictos y versionado.

No utilicé la IA para reemplazar la ejecución del trabajo. Los comandos y procedimientos fueron ejecutados y comprobados en mi propio entorno.

Cuando recibí indicaciones de la IA, las contrasté con el enunciado del TP y comprobé que los cambios realizados en el repositorio fueran los esperados.
