FROM ubuntu:18.04

WORKDIR /srv/node

# Install golang compiler/tools
RUN apt-get update && \
    apt-get install golang-go -y && \
    apt-get install iputils-ping -y

ENV GOPATH=/srv/node

EXPOSE 80

# Compile code
ADD src ./src
RUN go build -o ./build/out.o ./src/main.go

# Run binary output from compiler
CMD ["./build/out.o"]
