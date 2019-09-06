FROM ubuntu:18.04

WORKDIR /srv/node

# Install golang compiler/tools
RUN apt-get update
RUN apt-get install golang-go -y

ENV GOPATH=/srv/node

# Compile code
ADD src ./src
RUN go build -o ./build/out.o ./src/main.go

# Run binary output from compiler
CMD ["./build/out.o"]
