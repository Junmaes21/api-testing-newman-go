# API de libros: pruebas automatizadas con Newman

Proyecto didáctico para aprender a probar una API REST con una colección de Postman ejecutada desde Newman. La API se construye con Go y se despliega como servicio web público en Render.

[![Deploy to Render](https://render.com/images/deploy-to-render-button.svg)](https://dashboard.render.com/blueprint/new?repo=https%3A%2F%2Fgithub.com%2FJunmaes21%2Fapi-testing-newman-go)

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

1. Abre el botón **Deploy to Render** de arriba para cargar el Blueprint de este repositorio.
2. Si Render lo solicita, autoriza acceso de lectura al repositorio de GitHub y elige el workspace correcto.
3. Revisa que el servicio use el plan gratuito y selecciona **Apply** para crear el Blueprint. Render desplegará los commits cuando las verificaciones de GitHub Actions pasen.
4. La API desplegada está en [api-testing-newman-go.onrender.com](https://api-testing-newman-go.onrender.com).
5. Comprueba el endpoint `/health` y ejecuta la colección Newman contra la URL pública:

```bash
newman run postman/api-testing-newman-go.postman_collection.json \
  --env-var "baseUrl=https://api-testing-newman-go.onrender.com"
```

Render puede suspender servicios del plan gratuito cuando no reciben tráfico; la primera solicitud posterior puede tardar mientras vuelve a iniciar. El servicio es demostrativo y no debe usarse para guardar datos reales.

## Imágenes preparadas para el artículo

1. [Diagrama del flujo de pruebas](article/images/01-flujo-api.png)
2. [Árbol del repositorio público](article/images/02-project-github.png)
3. [Resumen de Newman, 8 solicitudes y 14 aserciones](article/images/03-newman-results.png)
4. [Ejecución aprobada en GitHub Actions](article/images/04-github-actions.png)
5. [Respuesta del endpoint público `/health`](article/images/05-api-health.png)

El artículo indica el punto exacto para subir cada archivo a Medium.

## Licencia

Este material se comparte con fines educativos.
