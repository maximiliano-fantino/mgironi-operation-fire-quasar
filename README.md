# operation-fire-quasar

El programa retorna la fuente y contenido del mensaje de auxilio. El mismo puede ser ejecutado en modo programa comando o en modo servidor web. En el caso del modo servidor web el servicio puede ser consumido via api rest (v2.0.0)

# stack tecnologico

El proyecto esta implementado en *Golang* 1.17

El framework web HTTP es *Gin (https://github.com/gin-gonic/gin)*

Las pruebas unitarias se implementaron con las bibliotecas *testing* y *net/http/httptest*

Tanto el entorno de desarrollo como ejecucion se uso *linux/unix*

Para el entorno local se puede usar (aparte de las tools propias de go) el Dockerfile para permitir el build de la imagen local y tambien docker-compose (en modo servidor web). La imagen de docker generada es similar la misma que se utiliza en el despliegue en entorno del proveedor cloud elegido. 

La solucion cloud que se uso para disponibilizar el servicio es *Google Run* de Google Cloud Platform. El servicio operation-fire-quasar esta conectado/configurado con el repositorio de github para hacer el build y deploy automatico, segun los eventos configurados. Ver el siguiente diagrama orientativo

<p align="center">
<img src="https://user-images.githubusercontent.com/40694446/151864237-12bb0fb8-32c0-4fbc-bb28-a0e3b4e3dbda.png"
</p>

# documentacion 

La docuemntacion de uso de la api rest (ejecucion en modo servidor web) esta documentada con swagger y la misma esta disponible en el mismo servicio desplegado. Es accesible a traves del path '/swagger/index.html'

La documenatcion de uso del programa comando (ejecucion en modo programa comando) se encuentra disponible como menu de ayuda del programa. Es accesible a traves del parametro '-h' o 'help'


# ejecucion en modo programa comando

El programa puede ser ejecutado en modo programa comando (luego de haber sido instalado). El mismo devuelve en consola el resultado de los calculos. 
 
Para realizar los calculos, el programa comando se puede ejecutar de la siguiente manera

    $ operation-fire-quasar -distances=100,200.65,-300.47 -message=this..the.complete.message,.is.the..message,.is...message

Para mas detalles de cada argumento, se recomeinda ejecutar el menu de ayuda con el siguiente comando

    $ operation-fire-quasar -h

# ejecucion en modo web server

El programa puede ser ejecutado en modo de servidor web. Se peude especificar el puerto de atencion de las peticiones http a traves de la variable de entorno 'PORT'
El modo se activa usando el argumento -profile=server, de la siguiente manera.

    $  operation-fire-quasar -profile=server

Los endpoints disponibles se encuentran disponibles via swagger en la ruta 

https://operation-fire-quasar-srv-lr7wlwx33q-ue.a.run.app/swagger/index.html

# parametrizacion de informacion de satelites

En ambos modos es posible parametrizar la infomacion de las ubicaciones de los distintos satelites a traves de las siguientes variables de entorno:

    . OFQ_KENOBI
    . OFQ_SKYWALKER
    . OFQ_SATO

El formato a utilizar en dichas variables es '<name>_<xcoord>,<ycoord>' . Ejemplo: 'kenobi_100.23,-287.15'
    
# administracion en google cloud platform

El servidor web se encuentra desplegado en el servicio Google Run. Se encuentra configuardo el build y despliegue automaticos usando como fuente el repositorio privado en github, segun los eventos configurados   

Tanto el despliegue y actualizacion del servicio se puede aplicar usando el archivo service.yaml el cual esta basado en el componenete service de kubernetes 
    
    $ gcloud run services replace environments/gcloud/service.yaml

La publicacion y accesibilidd al servicio se puede actualizar usando el archivo policy.yaml
    
    $ gcloud run services set-iam-policy operation-fire-quasar-srv environments/gcloud/policy.yaml

# pruebas

Se encuentran implementadas pruebas unitarias, la ejecucion de las mismas se realizan a traves de go test tool. Un resumen de la ejecucion en todo el proyecto y el coverage en cada package se puede acceder a traves del siguiente comando:

    $ go test ./...  -coverprofile=c.out

Tambien es posible ejecutar una prueba directa al servicio desplegado usando el comando curl o cualquier cliente apirest, tomando los archivos json de pruebas (los mismos son usados para las pruebas con la biblioteca 'net/http/httptest'). El siguiente es un ejemplo usando curl.

    $ curl -X POST -H "Content-Type: application/json" -d @_test/topSecret_test1_request.json https://operation-fire-quasar-srv-lr7wlwx33q-ue.a.run.app/topsecret/


# algunos comandos docker para uso en el ambiente local

. Construccion de imagen docker con tag 
    
    $ docker build . --tag operation-fire-quasar:1.0.0

. Ejecutacion de la app standalone
    
    $ docker run -it --rm con docker (-p 8080:3001 --name operation-fire-quasar-running operation-fire-quasar:1.0.0

# algunos comandos docker-compose para uso en el ambiente local

se debe acceder a 'environments/local'

    . Cosntruccion de la imagen docker
        
        $ docker-compose build

    . Inicializacion de la instancia
        
        $ docker-compose up

    . Desmontaje de la instancia
        
        $ docker-compose down
