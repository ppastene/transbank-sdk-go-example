# Transbank SDK Go Example

English version below

Este es un ejemplo de implementación del [SDK de Transbank](https://github.com/ppastene/transbank-sdk-go) escrito en Go. Hace uso de [Goravel](https://github.com/goravel/goravel) para el manejo de rutas, vistas y controladores.

## Requisitos
- Go 1.24.0
- [Air](https://github.com/air-verse/air) (opcional)
- Docker (opcional)

## Instalación
1. Haga un git clone de este proyecto o descarguelo.
2. Copie el archivo .env.example y renombrelo a .env.
3. En la raiz del proyecto ejecute el comando ```./artisan key:generate```. Ignore las advertencias acerca de la base de datos ya que este proyecto no hace uso de una.
4. Ejecute ```go mod tidy```
5. Ejecute ```go run .```. Se recomienda el uso de Air para la ejecución de este proyecto (ver abajo)
6. Acceda a la ruta http://localhost:3000

Se recomienda el uso de [Air](https://github.com/air-verse/air) para su uso en vez de usar ```go run .```. Siga las instrucciones de instalación, luego ejecute ```air init``` del cual creará un archivo ```.air.toml``` y finalmente ejecute el proyecto con ```air```

Este proyecto provee de una imagen Docker y de un archivo docker-compose. Si va a usar Docker cambie en el archivo .env el valor de ```APP_HOST``` por 0.0.0.0, luego ejecute ``` docker compose up```.

Consulte la [documentación de Goravel](https://www.goravel.dev/getting-started/installation.html) para mas información

-------------------------------------------------------------------

This is a [Transbank SDK](https://github.com/ppastene/transbank-sdk-go) implementation example written in the Go language. Make use of [Goravel](https://github.com/goravel/goravel) for route management, views and controllers.

## Requirements
- Go 1.24.0
- [Air](https://github.com/air-verse/air) (optional)
- Docker (optional)

## Installation
1. Make a git clone of this project or download it.
2. Copy the .env.example file and rename it to .env
3. In the root project execute the ```./artisan key:generate``` command. Ignore the database warnings since this project don't use one.
4. Execute ```go mod tidy```
5. Execute ```go run .```. It's recommendated the use of the Air package to run this project (see below)
6. Access to the route http://localhost:3000

It's recommended the use of the [Air](https://github.com/air-verse/air) package instead of ```go run .```. Follow the install instructions, then execute ```air init``` which it creates the ```.air.toml``` file, and finally run the project with ```air```

This project provides a Docker image and a docker-compose file. If you wanna use them in the .env file change the ```APP_HOST``` value with 0.0.0.0, then execute ```docker compose up```.

Follow the [Goravel documentation](https://www.goravel.dev/getting-started/installation.html) for more information