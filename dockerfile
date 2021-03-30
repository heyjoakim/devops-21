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

# Make port 80 available to the host
# Make port 2222 available for SSH from Azure
# Make port 21 available for FTPS to access logs 
EXPOSE 8000 2222 21

# Build and run the server when the container is started
RUN go build /src/minitwit.go
ENTRYPOINT ./minitwit
