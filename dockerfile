#image base Golang
FROM golang:1.20.4

# Crear un usuario sin privilegios
RUN groupadd -r appgroup && useradd -r -g appgroup appuser

# Especificar los recursos que necesitamos
ENV CPU_SET 0
ENV MEMORY 512m
ENV DEVICE /dev/sda
ENV PID host
ENV USER appuser

# Dar acceso a los dispositivos GPU
--gpus all

WORKDIR /GUI
#Lits moduls usage
COPY src/go.mod .
COPY src/go.sum .

#Download moduls
RUN go mod download

#Elements from conteiner
COPY . .

EXPOSE 8080
#Command compile to binary
RUN go build -o GUI .
#Run binary
CMD [ "./GUI" ]

#docker run --cpuset-cpus $CPU_SET --memory $MEMORY --device $DEVICE --pid $PID --user $USER -p 8080:8080 my-image