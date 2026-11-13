# Probar una API REST con Newman y automatizar el flujo en GitHub Actions

## De una petición que “funciona” a una API que podemos volver a comprobar

Cuando empezamos a construir una API, es fácil comprobar que responde con una llamada desde Postman y dar el trabajo por terminado. Sin embargo, esa prueba manual solo nos cuenta qué ocurrió en ese momento. En el siguiente cambio alguien podría modificar una respuesta, romper una ruta o aceptar datos incompletos, y quizá nos enteremos cuando otra parte de la aplicación ya dependa de ese comportamiento.

Por eso, en este artículo vamos a recorrer un ejemplo pequeño pero completo: construiremos una API REST de libros en Go, describiremos sus escenarios en una colección de Postman, ejecutaremos esa colección con Newman y haremos que GitHub Actions repita las verificaciones automáticamente. Después dejaremos la aplicación lista para publicarse como servicio web en Render.

La idea no es cubrir cada detalle de una plataforma, sino mostrar cómo conectar las piezas para que una comprobación que hoy hacemos a mano también pueda ejecutarse mañana, después de cada cambio.

> **[IMAGEN 1: subir `images/01-flujo-api.png` aquí]**
> Diagrama del recorrido: aplicación Go → colección Postman → Newman → GitHub Actions → Render.

## ¿Por qué Newman para probar una API?

Hay varias herramientas que sirven para revisar APIs, y cada una aporta algo distinto. Cypress resulta muy útil cuando queremos probar interacciones de usuario en el navegador y también puede enviar solicitudes HTTP desde sus pruebas. k6 destaca cuando la pregunta principal es cómo responde un servicio bajo carga. En este ejemplo nos interesa describir escenarios funcionales de HTTP y poder reutilizar una colección de Postman tanto en local como en integración continua, así que Newman encaja de forma directa.

Newman es el ejecutor de línea de comandos de las colecciones de Postman. Esto permite guardar las solicitudes y sus comprobaciones como archivos del proyecto, versionarlas junto con el código y ejecutarlas en una terminal o en un pipeline. Así evitamos que el conocimiento de las pruebas quede solamente en una colección guardada en la computadora de alguien.

Nuestra API administrará libros. Tendrá un endpoint de salud y operaciones para listar, consultar, crear, actualizar y eliminar recursos. Aunque el dominio es sencillo, nos alcanza para comprobar códigos de estado, contenido JSON, datos de respuesta y un flujo completo que depende de varias solicitudes.

## El proyecto y la API

La estructura separa el punto de entrada HTTP, la lógica de la API, la colección de pruebas y el workflow de GitHub Actions. Esta separación ayuda a encontrar cada pieza rápido, sobre todo cuando volvamos al proyecto para agregar un escenario nuevo.

> **[IMAGEN 2: subir `images/02-project-github.png` aquí]**
> Captura real del árbol de archivos del repositorio público.

La API expone `/health` para indicar que el proceso está activo y agrupa los recursos de libros bajo `/api/books`. Para mantener el ejemplo enfocado, los libros se guardan en memoria. Al reiniciar el proceso, vuelven los datos iniciales, por lo que esta implementación es apropiada para aprender y demostrar pruebas, pero no para almacenar información real.

Podemos levantar el servidor desde la raíz del repositorio con Go:

```bash
go run ./cmd/api
```

La aplicación escucha en el puerto 8080 de forma predeterminada. En Render puede utilizar la variable `PORT` que proporciona la plataforma. Un primer chequeo manual nos confirma que el servicio está disponible:

```bash
curl http://localhost:8080/health
```

La respuesta esperada es un objeto JSON con el estado `ok`. A partir de aquí, la colección automatizará esta comprobación y otras que cubren la funcionalidad de libros.

## Convertir los requisitos en escenarios de API

Antes de escribir solicitudes, conviene preguntarnos qué comportamiento queremos proteger. Para este ejemplo vamos a verificar que el health check responde correctamente, que podemos listar y consultar libros, que crear un recurso devuelve su identificador, que la actualización conserva los datos enviados y que una eliminación deja de encontrar ese recurso. También comprobaremos que la API rechaza una solicitud cuando falta un campo obligatorio.

La colección de Postman guarda esas solicitudes en `postman/api-testing-newman-go.postman_collection.json`. Dentro de cada solicitud, los tests de Postman pueden inspeccionar tanto la respuesta como sus propiedades. Por ejemplo, después de crear un libro obtenemos el identificador generado y lo reutilizamos en las solicitudes de actualización y eliminación. De esa manera, los últimos pasos trabajan con el mismo recurso que creó la prueba, en lugar de depender de un identificador escrito a mano.

Un fragmento representativo de una comprobación valida el código HTTP y el contenido JSON:

```javascript
pm.test('responde 201', () => pm.response.to.have.status(201));

const book = pm.response.json().data;
pm.test('el libro creado tiene id y título', () => {
  pm.expect(book.id).to.be.a('number');
  pm.expect(book.title).to.eql('El mundo es ancho y ajeno');
});

pm.collectionVariables.set('createdBookId', book.id);
```

La aserción no se queda en “el servidor contestó”. Comprueba que la creación produjo el resultado esperado y prepara el dato que necesitan las siguientes llamadas. Ese pequeño detalle es útil porque convierte varias solicitudes sueltas en un flujo reproducible.

En local, podemos ejecutar la colección así:

```bash
npm install --global newman@6.2.1
newman run postman/api-testing-newman-go.postman_collection.json \
  --env-var "baseUrl=http://localhost:8080"
```

La dirección base se pasa como variable para que la misma colección funcione en la computadora de desarrollo, en GitHub Actions y contra el servicio publicado. Esto evita mantener copias de la colección con URLs distintas.

> **[IMAGEN 3: subir `images/03-newman-results.png` aquí]**
> Resumen de Newman: ocho solicitudes, catorce aserciones y cero fallos, según el log de GitHub Actions.

## Hacer que GitHub Actions repita las pruebas

Una vez que el flujo funciona localmente, el siguiente paso es pedirle a GitHub que lo ejecute en cada `push` y en las solicitudes de cambio a la rama principal. El workflow del proyecto está en `.github/workflows/ci.yml` y combina varias comprobaciones: revisa el formato del código Go, ejecuta pruebas unitarias con cobertura, compila la aplicación y, finalmente, lanza Newman.

La parte importante del orden es iniciar primero la API dentro del runner. El workflow espera hasta que `/health` responda y entonces invoca Newman con `baseUrl` apuntando al proceso local. Si una aserción falla, Newman termina con un código de salida distinto de cero, por lo que el job queda marcado como fallido y el equipo puede revisar el cambio antes de integrarlo.

Este fragmento muestra la llamada central del workflow:

```yaml
- name: Newman API tests
  run: >-
    newman run postman/api-testing-newman-go.postman_collection.json
    --env-var "baseUrl=http://127.0.0.1:8080"
```

En el repositorio también se configuran Go y Node.js, porque la aplicación se compila con Go y Newman se ejecuta sobre Node. Mantener estas versiones declaradas hace más predecible el entorno de CI y permite que otra persona pueda repetir el flujo siguiendo el README.

> **[IMAGEN 4: subir `images/04-github-actions.png` aquí]**
> Captura real del job de GitHub Actions con sus pasos completados.

La automatización no elimina la necesidad de revisar las pruebas. Si los requisitos de la API cambian, los escenarios también deberían actualizarse. La ventaja es que ahora el equipo recibe una señal inmediata cuando el comportamiento que decidimos cuidar deja de cumplirse.

## Publicar la API para que otras personas puedan probarla

El proyecto incluye un `Dockerfile` y un archivo `render.yaml` para crear un servicio web en Render. Después de publicar el repositorio en GitHub, abrimos el botón **Deploy to Render** del README, autorizamos el acceso de lectura al repositorio si Render lo solicita y revisamos la configuración del Blueprint. Al seleccionar **Apply**, Render crea el servicio. La configuración espera que las verificaciones de GitHub Actions pasen antes de desplegar cada commit; cuando el primer despliegue termina, la API queda disponible bajo una URL pública de `onrender.com`.

Cuando el despliegue termine, hay que abrir `/health` en esa dirección y confirmar que responde. Luego podemos ejecutar la misma colección contra el servicio remoto, cambiando solamente `baseUrl`:

```bash
newman run postman/api-testing-newman-go.postman_collection.json \
  --env-var "baseUrl=https://api-testing-newman-go.onrender.com"
```

> **[IMAGEN 5: subir `images/05-api-health.png` aquí]**
> Respuesta HTTP verificada del endpoint público de salud.

Un servicio gratuito puede suspenderse después de un periodo sin tráfico y tardar un poco más en responder a la primera solicitud posterior. Es una limitación que debemos explicar cuando compartimos el enlace, especialmente si alguien espera una respuesta instantánea. La instancia de este ejemplo está disponible en `https://api-testing-newman-go.onrender.com`.

## Qué conseguimos con este recorrido

Con un ejemplo acotado conectamos cuatro prácticas: una API que podemos iniciar localmente, escenarios funcionales guardados con el proyecto, una ejecución automatizada que informa si algo se rompe y un servicio accesible desde Internet. Además, los mismos escenarios sirven en los tres contextos, porque cambiamos la URL base en lugar de duplicar las pruebas.

El proyecto deja espacio para crecer. Podemos añadir autenticación, persistencia con una base de datos, pruebas de carga con k6 o validaciones de contrato más estrictas. Para cada extensión, lo importante será elegir qué riesgo queremos detectar y añadir una comprobación que pueda repetirse de forma confiable.

Si quieres reproducir el ejemplo, revisa primero el README, ejecuta la API, corre Newman y después observa el job de GitHub Actions. Una vez desplegado el servicio, cambia la URL base y vuelve a lanzar la colección. Así podrás ver el mismo flujo desde tu máquina hasta un endpoint público.

## Recursos

- Código del proyecto: [github.com/Junmaes21/api-testing-newman-go](https://github.com/Junmaes21/api-testing-newman-go)
- API desplegada: [api-testing-newman-go.onrender.com](https://api-testing-newman-go.onrender.com)
- [Documentación oficial de Newman](https://learning.postman.com/docs/collections/using-newman-cli/newman-options)
- [Documentación oficial de GitHub Actions](https://docs.github.com/actions)
- [Guía oficial de Render para servicios web Go](https://render.com/docs/deploy-go-nethttp)
