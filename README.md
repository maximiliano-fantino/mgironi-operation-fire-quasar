# operation-fire-quasar

El programa retorna la fuente y contenido del mensaje de auxilio. El mismo puede ser ejecutado en modo programa comando o en modo servidor web. En el caso del modo servidor web el servicio puede ser consumido via api rest.

última versión v3.0.0

# stack tecnologico

. El proyecto esta implementado en *Golang* 1.17.

. El framework web HTTP usado es *Gin (https://github.com/gin-gonic/gin)*.

. Para almacenamiento temporal de las peticiones por partes (split) se utiliza Redis

. Las pruebas unitarias se implementaron con las bibliotecas *testing*, *net/http/httptest* y *redigomock* según se necesita en cada caso.

. Tanto para el entorno de desarrollo como de despliegue se usa *linux/unix*.

. Para el entorno local se puede usar (aparte de las tools propias de go) el Dockerfile para permitir el build de la imagen local y tambien docker-compose (en modo servidor web). La imagen de docker generada es similar la misma que se utiliza en el despliegue en entorno del proveedor cloud elegido. 

. La solucion cloud que se uso para disponibilizar el servicio es *Google Run* de Google Cloud Platform, como también servicios adicionales para facilitar el build y el despliegue del mismo proveedor. El servicio operation-fire-quasar esta conectado/configurado con el repositorio de github para hacer el build y deploy automatico, segun los eventos configurados. Ver el siguiente diagrama orientativo

<p align="center">
<img src="https://user-images.githubusercontent.com/40694446/151864237-12bb0fb8-32c0-4fbc-bb28-a0e3b4e3dbda.png"
</p>

# documentacion 

Aparte del presente documento, se cuenta con documentación propia según el perfil de uso del programa, y es la siguiente.
 
La docuemntacion de uso de la api rest (ejecucion en modo servidor web) esta documentada con swagger y la misma esta disponible en el mismo servicio desplegado. Es accesible a traves del path '/swagger/index.html'

La documenatcion de uso del programa comando (ejecucion en modo programa comando) se encuentra disponible como menu de ayuda del programa. Es accesible a traves del parametro '-h' o 'help'

# principales funciones del programa
 
## calculo para determinación de la ubicación
 
Para poder determinar la ubicación y considerando que se cuenta con las distancias a tres puntos cuyas coordenadas son conocidas se aplica el método matemático de trilateración. El mismo se describe como la intersección de tres esferas con centro en los puntos conocidos y de radios la distancia a cada uno. En este caso particularmente solo se cuentan con dos dimensiones con lo que en lugar de esferas se opera con circunferencias. Adicionalmente, los puntos conocidos no se encuentran alineados en un mismo eje (almenos dos) por la tanto es necesario realizar una rotación de los ejes (aparte de la traslación que propone propia el método). De esta forma se puede prevenir un error de cálculo
 
## armado del mensaje emitido

El mesnaje emitido,el cual es recibido en partes (una por cada satelite) se trata de la siguiente manera:
1. Se verifica cual de los mesnajes tiene mas palabras. Dado que se reciben areglos de strings y si una palabra es faltante se mantiene la posicion.
2. En caso de encontrar que alguno de los mensajes es mas corto se (tiene menos palabras) se lo completa (exclusivamente al inicio) con nuevas posciones de cadena vacia.
3. Se realiza una iteracion tomando las palabras de cada mensaje, las que no estan vacias y verificando que las palabras en cada posicion no sean distintas (exceptuando la cadena vacia)
4. Se toma el arreglo resultante como arreglo dle mensaje completo. 

## tratamiento de llamadas por partes *split*

Para recibir los datos de cada satelite por separado se dispone de los siguientes dos endpoints. Uno recibe la informacion de los satelites y el otro permite obtener el resultado del calculo que previamente fue recolectado de a partes. Y son los siguientes:

. POST /topsecret_split/\[operation\]
. GET /topsecret_split/{operation}

### POST /topsecret_split/\[operation\]

Recibe el dato de un satelite y retorna el codigo de operacion para poder seguir enviando datos de los resto de los satelites. El paramtero \[operation\] es opcional dado que se dispone de un mecanismo deteccion a traves de la porcion del mensaje en caso de que no se cuente con el codigo de operacion. Asi mismo, al enviar el primer dato de un satelite el codigo de operacion se genera y es devuelto.
Es prefrible la utilizacion de codigo de operacion para los subsiguientes envios de datos hasta completar el total de satelites requerido.

### GET /topsecret_split/{operation}

Retorna los resultados del calculo siempre que se hallan completado los datos de los satelites.

A continuacion se presenta un diagrama de actividad para resumir la combinacion de ambos endpoints 

![topsecret-topsecret_split](https://user-images.githubusercontent.com/40694446/152417245-3c776296-d694-4808-82ea-61126ee4291c.png)

Para mas detalles sobre las llamadas a la api ver .. https://operation-fire-quasar-srv-lr7wlwx33q-ue.a.run.app/swagger/index.html 

---
 
# ejecucion en modo programa comando

El programa puede ser ejecutado en modo programa comando (luego de haber sido instalado), o en su defecto con go run. El mismo devuelve en consola el resultado de los calculos. 
 
Para realizar los calculos, el programa comando se puede ejecutar de la siguiente manera

    $ operation-fire-quasar -distances=100,200.65,-300.47 -message=this..the.complete.message,.is.the..message,.is...message

Para mas detalles de cada argumento, se recomeinda ejecutar el menu de ayuda con el siguiente comando

    $ operation-fire-quasar -h

# ejecucion en modo web server

El programa puede ser ejecutado en modo de servidor web. Se peude especificar el puerto de atencion de las peticiones http a traves de la variable de entorno 'PORT'
El modo se activa al agregar el argumento -profile=server, de la siguiente manera.

    $  operation-fire-quasar -profile=server

Los endpoints disponibles se encuentran documentados y disponibles para ser probados via swagger en la siguiente ruta 

https://operation-fire-quasar-srv-lr7wlwx33q-ue.a.run.app/swagger/index.html

# parametrizacion de informacion de satelites

En ambos modos, es posible parametrizar la infomacion de las ubicaciones de los distintos satelites a traves de las siguientes variables de entorno:

    . OFQ_KENOBI
    . OFQ_SKYWALKER
    . OFQ_SATO

El formato a utilizar en dichas variables es *name>_xcoord,ycoord* . Ejemplo: *kenobi_100.23,-287.15*
    
# administracion en google cloud platform

El servidor web se encuentra desplegado en el servicio Google Run. Se encuentra configuardo el build y despliegue automaticos, usando como fuente el repositorio privado en github. Dichas operacion se inician segun los eventos configurados.

Tanto el despliegue y actualizacion del servicio se puede aplicar usando el archivo service.yaml el cual esta basado en el componenete service de kubernetes 
    
    $ gcloud run services replace environments/gcloud/service.yaml

La publicacion y accesibilidd al servicio se puede actualizar usando el archivo policy.yaml
    
    $ gcloud run services set-iam-policy operation-fire-quasar-srv environments/gcloud/policy.yaml

# pruebas

Se encuentran implementadas pruebas unitarias, la ejecucion de las mismas se realizan a traves de go test tool. Para un resumen de la ejecucion en todo el proyecto y el coverage en cada package se puede acceder a traves del siguiente comando:

    $ go test ./...  -coverprofile=c.out

Tambien es posible ejecutar una prueba directa al servicio desplegado usando el comando curl o cualquier cliente apirest, tomando los archivos json de pruebas (los mismos son usados para las pruebas con la biblioteca 'net/http/httptest'). El siguiente es un ejemplo usando curl.

    $ curl -X POST -H "Content-Type: application/json" -d @_test/topSecret_test1_request.json https://operation-fire-quasar-srv-lr7wlwx33q-ue.a.run.app/topsecret/

----

# algunos comandos docker para uso en el ambiente local

## Construccion de imagen docker con tag 
    
    $ docker build . --tag operation-fire-quasar:1.0.0

## Ejecutacion de la app standalone
    
    $ docker run -it --rm con docker (-p 8080:3001 --name operation-fire-quasar-running operation-fire-quasar:1.0.0

# algunos comandos docker-compose para uso en el ambiente local

se debe acceder a 'environments/local'

## Cosntruccion de la imagen docker
        
    $ docker-compose build

## Inicializacion de la instancia
        
    $ docker-compose up

## Desmontaje de la instancia
        
    $ docker-compose down

## starts a redis-server (standalone)
    
    $ docker run --name some-redis -d redis:6.0-alpine redis-server --save 60 1 --loglevel warning

## Connect to local redis (standalone), to use redis-cli from console
    
    $ docker network create redis-ntk
    $ docker network connect redis-ntk some-redis
    $ docker run -it --network redis-ntk --rm redis:6.0-alpine redis-cli -h some-redis

## clean docker volumes
    
    $ docker volume ls
    $ docker volume rm VOLUME-NAME

## Starts app and redis with docker-compose
    
    $ docker-compose up
    $ docker-compose down
