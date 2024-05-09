# Start from a base image with Chrome installed
FROM selenium/standalone-chrome:latest

# Install Go
RUN sudo apt-get update && sudo apt-get install -y wget
RUN wget https://dl.google.com/go/go1.18.linux-amd64.tar.gz
RUN sudo tar -xvf go1.18.linux-amd64.tar.gz
RUN sudo mv go /usr/local

# Set environment path
ENV GOROOT=/usr/local/go
ENV GOPATH=$HOME/go
ENV PATH=$GOPATH/bin:$GOROOT/bin:$PATH

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.* ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o go-get-site-cli .

# Command to run the executable
CMD ["./go-get-site-cli"]
