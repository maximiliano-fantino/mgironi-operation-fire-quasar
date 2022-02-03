# operation-fire-quasar

El programa retorna la fuente y contenido del mensaje de auxilio. El mismo puede ser ejecutado en modo programa comando o en modo servidor web. En el caso del modo servidor web el servicio puede ser consumido vía api rest.

URL: https://operation-fire-quasar-srv-lr7wlwx33q-ue.a.run.app/swagger/index.html

versión actual v3.0.0

_(con el objeto de facilitar la evaluacion, el acceso al servidor es público)_

# stack tecnológico

. El proyecto esta implementado en *Golang* 1.17.

. El framework web HTTP usado es *Gin (https://github.com/gin-gonic/gin)*.

. Se utiliza *Redis* para almacenamiento temporal para el tratamiento de las peticiones por partes (split).

. Las pruebas unitarias se implementaron con las bibliotecas *testing*, *net/http/httptest* y *redigomock* según se necesita en cada caso.

. Tanto para el entorno de desarrollo como el de despliegue se uso *linux/unix*.

Para el entorno local se puede usar (aparte de las tools propias de go) el Dockerfile para permitir el build de la imagen local y también docker-compose (en modo servidor web). La imagen de docker generada es similar la misma que se utiliza en el despliegue en entorno del proveedor cloud elegido. 

La solucion cloud que se uso para disponibilizar el servicio es *Google Run* de Google Cloud Platform, también se usan servicios adicionales para facilitar el build y el despliegue, todo del mismo proveedor. El servicio operation-fire-quasar esta conectado/configurado con el repositorio privado de github para hacer el build y deploy automático, según los eventos configurados. Ver el siguiente diagrama orientativo

<p align="center">
<img src="https://user-images.githubusercontent.com/40694446/151864237-12bb0fb8-32c0-4fbc-bb28-a0e3b4e3dbda.png"
</p>

# documentación

Aparte del presente documento, se cuenta con documentación según el perfil de uso del programa, y es la siguiente.
 
La docuemntacion de uso de la api rest (ejecución en modo servidor web) esta documentada con swagger y la misma esta disponible en el mismo servicio desplegado. Es accesible a través del path '/swagger/index.html'

La documenatción de uso del programa comando (ejecución en modo programa comando) se encuentra disponible como menu de ayuda del programa. Es accesible a traves del parametro '-h' o 'help'

# principales funciones del programa
 
## cálculo para determinación de la ubicación
 
Para poder determinar la ubicación y considerando que se cuenta con las distancias a tres puntos cuyas coordenadas son conocidas se aplica el método matemático de trilateración. El mismo se describe como la intersección de tres esferas con centro en los puntos conocidos y de radios de la distancia a cada uno de ellos. En este caso particular, sólo se cuenta con dos dimensiones con lo que en lugar de esferas se opera con circunferencias. Adicionalmente, como los puntos conocidos no se encuentran alineados en un mismo eje (almenos dos), es necesario realizar una rotación de los ejes (aparte de la traslación que propone el método en si). De esta forma se puede prevenir un error significativo en el cálculo.
 
## armado del mensaje emitido

El mesnaje emitido, el cual es recibido en partes (una por cada satelite) se trata de la siguiente manera:
1. Se verifica cual de los mesnajes tiene mas palabras. Dado que se reciben arreglos de strings y si una palabra es faltante se mantiene la posición.
2. En caso de encontrar que alguno de los mensajes es mas corto (tiene menos palabras) se lo completa (exclusivamente al inicio) con nuevas posciones de cadena vacía.
3. Se realiza una iteración tomando las palabras de cada mensaje, las que no estan vacias y verificando que las palabras en cada posicion no sean distintas (exceptuando la cadena vacía)
4. Se toma el arreglo resultante como arreglo del mensaje completo. 

## tratamiento de llamadas por partes *split*

Para recibir los datos de cada satelite por separado se dispone de dos endpoints. Uno recibe la información de los satélites y el otro permite obtener el resultado del cálculo que previamente fue recolectado de a partes. Y son los siguientes:

. POST /topsecret_split/\[operation\]
. GET /topsecret_split/{operation}

### POST /topsecret_split/\[operation\]

Recibe el dato de un satelite y retorna el codigo de operacion para poder seguir enviando datos de los resto de los satélites. El parámetro \[operation\] es opcional dado que se dispone de un mecanismo detección por aproximación de frases similares a través de la porción del mensaje en caso de que no se cuente con el código de operación. Asi mismo, al enviar el primer dato de un satélite el código de operación se genera y es retornado en la llamada.

Es recomendable la utilización del código de operación para los subsiguientes envios de datos hasta completar el total de satélites requerido, dado que por un lado el mécanismo de detección es por aproximación con lo que en un mensaje de pocas palabras puede inferir en un error de matcheo. Por otro lado el costo computacional es mucho mayor.

### GET /topsecret_split/{operation}

Retorna los resultados del cálculo siempre que se hubiera completado la recolección de los datos de los satélites.

A continuacion se presenta un diagrama de actividad para resumir la combinación de ambos endpoints.

![topsecret-topsecret_split](https://user-images.githubusercontent.com/40694446/152417245-3c776296-d694-4808-82ea-61126ee4291c.png)

Para mas detalles sobre las llamadas a la api, ver .. https://operation-fire-quasar-srv-lr7wlwx33q-ue.a.run.app/swagger/index.html 

---
 
# ejecución en modo programa comando

El programa puede ser ejecutado en modo programa comando (luego de haber sido instalado), o en su defecto con go run. El mismo devuelve en consola el resultado de los cálculos. 
 
Para realizar los cálculos, el programa comando se puede ejecutar de la siguiente manera

    $ operation-fire-quasar -distances=100,200.65,-300.47 -message=this..the.complete.message,.is.the..message,.is...message

Para mas detalles de cada argumento, se recomeinda ejecutar el menu de ayuda con el siguiente comando

    $ operation-fire-quasar -h

# ejecución en modo web server

El programa puede ser ejecutado en modo de servidor web. Se peude especificar el puerto de atención de las peticiones http a través de la variable de entorno 'PORT'.

 El modo se activa al agregar el argumento -profile=server, de la siguiente manera.

    $  operation-fire-quasar -profile=server

Los endpoints disponibles se encuentran documentados y listos para ser probados vía swagger en la siguiente ruta 

https://operation-fire-quasar-srv-lr7wlwx33q-ue.a.run.app/swagger/index.html

# parametrización de información de satélites

En ambos modos, es posible parametrizar la infomación de las ubicaciones de los distintos satélites a través de las siguientes variables de entorno:

    . OFQ_KENOBI
    . OFQ_SKYWALKER
    . OFQ_SATO

El formato a utilizar en dichas variables es *name>_xcoord,ycoord* . Ejemplo: *kenobi_100.23,-287.15*
    
# administración en google cloud platform

El servidor web se encuentra desplegado en el servicio Google Run. Y configurado el build y despliegue automáticos, se usa como fuente el repositorio privado en github. Dichas operaciones se inician según los eventos configurados. El servicio de google run cuenta con la capacidad de autoescalamiento y solo se consume computo al momento de atender las llamadas.

Tanto el despliegue y actualización del servicio se puede aplicar usando el archivo service.yaml de la siguiente forma.
    
    $ gcloud run services replace environments/gcloud/service.yaml

La publicación y accesibilidd al servicio se puede actualizar usando el archivo policy.yaml
    
    $ gcloud run services set-iam-policy operation-fire-quasar-srv environments/gcloud/policy.yaml

# pruebas

Se encuentran implementadas pruebas unitarias, la ejecución de las mismas se realizan a través de go test tool. Para un resumen de la ejecución en todo el proyecto y el coverage en cada package se puede acceder a través del siguiente comando:

    $ go test ./...  -coverprofile=c.out

Tambien es posible ejecutar una prueba directa al servicio desplegado usando el comando curl o cualquier cliente apirest, tomando los archivos json de pruebas (los mismos son usados para las pruebas con la biblioteca 'net/http/httptest'). El siguiente es un ejemplo usando curl.

    $ curl -X POST -H "Content-Type: application/json" -d @_test/topSecret_test1_request.json https://operation-fire-quasar-srv-lr7wlwx33q-ue.a.run.app/topsecret/

----

# algunos comandos docker para uso en el ambiente local

## Construcción de imagen docker con tag 
    
    $ docker build . --tag operation-fire-quasar:1.0.0

## Ejecutacion de la app standalone
    
    $ docker run -it --rm con docker (-p 8080:3001 --name operation-fire-quasar-running operation-fire-quasar:1.0.0

# algunos comandos docker-compose para uso en el ambiente local

se debe acceder a 'environments/local'

## Cosntrucción de la imagen docker
        
    $ docker-compose build

## Inicialización de la instancia
        
    $ docker-compose up

## Desmontaje de la instancia
        
    $ docker-compose down

## inicio de redis-server (standalone)
    
    $ docker run --name some-redis -d redis:6.0-alpine redis-server --save 60 1 --loglevel warning

## Conección a redis local (standalone), para redis-cli desde consola
    
    $ docker network create redis-ntk
    $ docker network connect redis-ntk some-redis
    $ docker run -it --network redis-ntk --rm redis:6.0-alpine redis-cli -h some-redis

## limpiar docker volumenes
    
    $ docker volume ls
    $ docker volume rm VOLUME-NAME

## Iniciar/detener app y redis con docker-compose
    
    $ docker-compose up
    $ docker-compose down
