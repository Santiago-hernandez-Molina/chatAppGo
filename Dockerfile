# Utiliza una imagen base de Go
FROM golang:latest

# Establece el directorio de trabajo dentro del contenedor
WORKDIR /app

# Copia los archivos de la aplicación al directorio de trabajo
COPY . .

# Instala las dependencias de la aplicación (si es necesario)
# RUN go get -d -v ./...

# Compila la aplicación
RUN go build -o main cmd/main.go

# Exponer el puerto en el que la aplicación se ejecutará
EXPOSE 8080

# Establecer las variables de entorno para la conexión a MongoDB

# Comando para ejecutar la aplicación cuando se inicie el contenedor
CMD ["./main"]
