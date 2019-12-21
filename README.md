# Proyecto Semestral
## Asignatura de Computación Paralela y Distribuida

### Segundo Semestre 2019 - Universidad Tecnológica Metropolitana

### Resumen del enunciado

Se solicita desarrollar un sistema distribuido, que permita encontrar el recorrido de micro que tenga menos combinaciones entre dos direcciones.

### Explicación de este repositorio

Este repositorio contiene el código fuente de un backend construido en Go, utilizando el framework echo.
La información de entregada fue preprocesada para ser utilizada en memoria en la ejecución del programa.
La búsqueda de los recorridos en micro se realizó con diferentes ciclos de búsqueda, automatizando su tiempo con Go routines.

##### Integrantes

* Valentina Faure
* Jacob Romero
* Pedro Valderrama
* Diego Sepúlveda

# Construir imagen de docker a partir de Dockerfile
docker build -t backend-proyecto-cpd:1.0 -f Dockerfile .

# Ejecutar contenedor en modo deamon exponiendo el puerto 80
docker run -d -p 80:8080 backend-proyecto-cpd:1.0
