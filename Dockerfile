# Use the offical Golang image as a build stage
FROM golang:1.21 AS build

# Set the working directory inside the container
WORKDIR /backend

#Copy go.mod and go.sum files
COPY go.mod go.sum ./

#Download dependencies
RUN go mod download

#Copy the source code in the container
COPY . .

# Build the Go app
RUN go build -o main ./backend/main.go # Build the go app from the main.go file

#Use a smaller image as a runtime stage
FROM alphine:latest

#Set the working directory in the container
WORKDIR /backend

# Copy the built app from the build stage
COPY --from=build /backend/main .

#Expose the port 8080
EXPOSE 8080

#Command to run the app
CMD ["./main"]