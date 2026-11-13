# API de libros: pruebas automatizadas con Newman

Proyecto didáctico para aprender a probar una API REST con una colección de Postman ejecutada desde Newman. La API se construye con Go y se despliega como servicio web público en Render.

[![Deploy to Render](https://render.com/images/deploy-to-render-button.svg)](https://render.com/deploy?repo=https://github.com/junmaes21/api-testing-newman-go)

## Qué incluye

- API REST de libros con operaciones para listar, consultar, crear, actualizar y eliminar.
- Pruebas unitarias con `go test`.
- Pruebas funcionales HTTP con Newman y una colección Postman versionada.
- GitHub Actions para validar formato, ejecutar pruebas, compilar la aplicación y correr Newman en cada `push` y `pull request` a `main`.
- `Dockerfile` y `render.yaml` para desplegar el servicio en Render.

## Requisitos

- Go 1.26 o posterior.
- Node.js 22 o posterior y Newman 6.2.1 para ejecutar la colección.

## Ejecutar la API

```bash
go run ./cmd/api
```

La API queda disponible en `http://localhost:8080`. Puedes cambiar el puerto con la variable de entorno `PORT`.

## Ejecutar las pruebas

```bash
go test -race -cover ./...
```

Con la API en marcha, ejecuta Newman:

```bash
npm install --global newman@6.2.1
newman run postman/api-testing-newman-go.postman_collection.json \
  --env-var "baseUrl=http://localhost:8080"
```

La colección cubre el health check, lectura, flujo CRUD, validación de datos y códigos HTTP. En GitHub Actions se usa el mismo comando contra la aplicación que inicia el propio workflow.

## Rutas principales

| Método | Ruta | Resultado esperado |
| --- | --- | --- |
| GET | `/health` | Estado del servicio |
| GET | `/api/books` | Lista de libros |
| GET | `/api/books/{id}` | Un libro o `404` |
| POST | `/api/books` | Crea un libro o responde `400` |
| PUT | `/api/books/{id}` | Actualiza un libro |
| DELETE | `/api/books/{id}` | Elimina un libro y responde `204` |

La aplicación guarda los datos en memoria para mantener el ejemplo pequeño. Al reiniciarse, se restauran los libros de demostración.

## Despliegue público en Render

1. Publica este repositorio en GitHub.
2. Inicia sesión en Render y elige **Deploy to Render** o crea un **Blueprint** desde este repositorio público.
3. Revisa el servicio definido en `render.yaml` y crea el despliegue.
4. Al finalizar el despliegue, usa la URL pública de Render, por ejemplo `https://<nombre-del-servicio>.onrender.com`.
5. Comprueba `https://<nombre-del-servicio>.onrender.com/health` y ejecuta la colección Newman indicando esa URL:

```bash
newman run postman/api-testing-newman-go.postman_collection.json \
  --env-var "baseUrl=https://<nombre-del-servicio>.onrender.com"
```

Render puede suspender servicios del plan gratuito cuando no reciben tráfico; la primera solicitud posterior puede tardar mientras vuelve a iniciar. El servicio es demostrativo y no debe usarse para guardar datos reales.

## Capturas sugeridas para el artículo

1. **Después de explicar la estructura del proyecto:** inserta una captura del árbol de carpetas mostrando `cmd`, `internal`, `postman` y `.github/workflows`.
2. **Después de presentar Newman:** inserta una captura del resultado de la colección con las ocho solicitudes aprobadas.
3. **Después de explicar CI:** inserta una captura del job `API checks` completado en GitHub Actions.
4. **En la sección de despliegue:** inserta una captura de la página `/health` abierta desde el dominio de Render.

## Licencia

Este material se comparte con fines educativos.
