FROM golang
ENV GO111MODULE=on

# Install any needed dependencies...
# RUN go get ...
# dependency using go mod
# Set the working directory
WORKDIR /src

COPY go.mod /src/.
COPY go.sum /src/.
RUN go mod download


# Copy the server code into the container

COPY . /src/.

# Make port 8080 available to the host
EXPOSE 80

# Build and run the server when the container is started
RUN go build /src/minitwit.go
ENTRYPOINT ./minitwit
