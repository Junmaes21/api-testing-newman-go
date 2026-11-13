# API de libros: pruebas automatizadas con Newman

Proyecto educativo para explorar una API REST de libros, probar sus operaciones con Newman y automatizar las validaciones con GitHub Actions.

## Enlaces

- **Aplicación publicada:** [api-testing-newman-go.onrender.com](https://api-testing-newman-go.onrender.com)
- **Artículo en Medium:** [De una API a producción: pruebas automatizadas con Newman](https://medium.com/@juniormaes21/de-una-api-a-producci%C3%B3n-pruebas-automatizadas-con-newman-d7fd498f617b)

La página principal permite consultar, buscar, agregar, editar y eliminar libros. La API también puede probarse con Postman o Newman. Los datos se guardan en memoria y vuelven a los valores de demostración cuando el servicio se reinicia.

## Tecnologías y automatización

- API e interfaz web implementadas en Go.
- Colección de Postman ejecutada desde Newman para verificar los endpoints y el flujo CRUD.
- GitHub Actions ejecuta pruebas, formato, compilación y Newman en cada `push` y `pull request` a `main`.
- Render publica la aplicación mediante la configuración de `render.yaml`.

## Ejecutar localmente

Requisitos: Go 1.26 o posterior, Node.js 22 o posterior y Newman 6.2.1.

Inicia la aplicación desde la raíz del repositorio:

```bash
go run ./cmd/api
```

Abre `http://localhost:8080`. Para correr las pruebas de Go:

```bash
go test -race -cover ./...
```

Con la aplicación en marcha, instala y ejecuta Newman:

```bash
npm install --global newman@6.2.1
newman run postman/api-testing-newman-go.postman_collection.json \
  --env-var "baseUrl=http://localhost:8080"
```

Para apuntar las pruebas a la aplicación publicada, cambia `baseUrl` por `https://api-testing-newman-go.onrender.com`.

## Endpoints principales

| Método | Ruta | Acción |
| --- | --- | --- |
| GET | `/health` | Comprueba el estado del servicio |
| GET | `/api/books` | Lista los libros |
| GET | `/api/books/{id}` | Consulta un libro |
| POST | `/api/books` | Crea un libro |
| PUT | `/api/books/{id}` | Actualiza un libro |
| DELETE | `/api/books/{id}` | Elimina un libro |

El servicio de Render puede tardar en responder cuando vuelve de un periodo sin tráfico. Esta aplicación es demostrativa y no debe usarse para almacenar datos reales.
