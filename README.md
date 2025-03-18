# Timeline Service
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=uala-challenge_timeline-service&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=uala-challenge_timeline-service)
![technology Go](https://img.shields.io/badge/technology-go-blue.svg)
![Viper](https://img.shields.io/badge/configuration-viper-green.svg)


## Descripción
El servicio **Timeline Service** es responsable de la gestión de las líneas de tiempo de los usuarios. Cada vez que un usuario sigue a otro, su línea de tiempo se actualiza. Para optimizar el rendimiento, este servicio consulta primero los IDs de los tweets almacenados en **Redis** y luego obtiene la información completa de cada tweet desde **DynamoDB**.

## Características principales
- Consulta rápida de timelines desde **Redis**.
- Obtención de detalles de los tweets desde **DynamoDB**.
- Actualización en tiempo real al seguir o dejar de seguir usuarios.
- API REST para consultar la línea de tiempo de los usuarios.
- Logs detallados para monitoreo y depuración.

## Tecnologías utilizadas
- **Golang**: Lenguaje de desarrollo principal.
- **Redis**: Almacenamiento en caché para IDs de tweets.
- **DynamoDB**: Base de datos para almacenar tweets completos.
- **Docker**: Contenedorización para despliegue.
- **REST API**: Exposición de endpoints.
- **Logrus**: Manejo de logs estructurados.

## Instalación y configuración
### Prerrequisitos
- Tener instalado **Go** (versión 1.18 o superior).
- Contar con **Docker** y **Docker Compose** (para pruebas locales).
- Configurar las variables de entorno adecuadas.

### Clonar el repositorio
```bash
  git clone https://github.com/uala-challenge/timeline-service.git
  cd timeline-service
```

### Configuración
El servicio utiliza un archivo de configuración en formato YAML. Antes de ejecutar, asegúrate de configurar correctamente `config.yaml`:
```yaml
aws:
  region: 'us-east-1'

dynamo:
  endpoint: 'http://localhost:4566'

redis:
  host: 'localhost'
  port: 6379
  db: 0
  timeout: 5
```

### Ejecutar en entorno local
#### Usando Go directamente
```bash
  go run main.go
```

#### Usando Docker
```bash
  docker-compose up --build
```

## API REST
El servicio expone los siguientes endpoints:

### **Consultar línea de tiempo de un usuario**
```
GET /timeline/{user_id}
```
#### **Descripción:**
Obtiene la línea de tiempo de un usuario, consultando primero en **Redis** y luego en **DynamoDB** para obtener los detalles de los tweets.

#### **Respuesta:**
```json
{"user_id":"user:354344","tweet_id":"tweet:9fa355fd-1e03-48d6-b923-1f985e9ca914","created":1741963493,"content":"Lorem ipsum dolor sit amet consectetur, adipisicing elit. Eveniet a necessitatibus voluptatum aliquam quasi iure et animi similique consequatur vero? Eum tenetur consequatur tempora exercitationem! Molestias, sapiente earum, qui ducimus et repellendus natus labore beatae archite."}
```

### **Actualizar línea de tiempo cuando un usuario sigue a otro**
```
PATCH /timeline/{user_id}
```
#### **Descripción:**
Este endpoint es llamado por **Users Service** cuando un usuario sigue a otro. Se actualiza la línea de tiempo agregando los últimos tweets del nuevo seguido.

#### **Payload esperado:**
```json
{
  "followed_user_id": "67890"
}
```

## Testing
Para ejecutar las pruebas unitarias:
```bash
  go test ./...
```
