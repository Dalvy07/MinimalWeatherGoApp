# syntax=docker/dockerfile:1.4-labs
# First stage - application build
FROM golang:alpine AS builder

# Set working directory
WORKDIR /app

# Install openssh-client and git for SSH access to GitHub
RUN apk add --no-cache git openssh-client

RUN mkdir -p /root/.ssh && \
    chmod 700 /root/.ssh && \
    ssh-keyscan github.com >> /root/.ssh/known_hosts

# THIS VERY IMPORTANT FOR OPTIMIZING THE SIZE OF THE IMAGE.
# When compiling with this setting as a result I get a statically linked binary, 
# which allows me to run this application using only a FROM scratch base image.
ENV CGO_ENABLED=0

# Variable for cache busting
# This is needed to force Docker to rebuild the image when the code changes
ARG CACHE_BUST=unknown
ARG BRANCH=main

# Сopy source files from GitHub repo
# Tryed using echo "Current time: $(date) because without it always use cache even if code was changed
# And I don't want to use --no-cache.
# But it still doesnt work
RUN --mount=type=ssh,id=github \
    echo "Cache bust: ${CACHE_BUST}" && \
    git clone --depth=1 --branch=${BRANCH} git@github.com:Dalvy07/MinimalWeatherGoApp.git .

# # Copy files from local folder
# COPY static/ ./static/
# COPY main.go go.mod ./ 

RUN --mount=type=secret,id=api_key \
    export API_KEY=$(cat /run/secrets/api_key) && \
    # Compile the application, flags for maximum size optimization
    go build -ldflags="-s -w -X main.apiKey=$API_KEY" -o weather-app .



# Second stage - creating minimal image
FROM scratch

# Define build argument with default value
ARG PORT=3000

# Set environment variables for port and author
ENV PORT=${PORT}
ENV APP_AUTHOR="Vladyslav Liulka <vladlulka3@gmail.com>"

# Set labels for the image metadata
LABEL org.opencontainers.image.authors="${APP_AUTHOR}"
LABEL org.opencontainers.image.title="Weather App"
LABEL org.opencontainers.image.description="A minimalistic weather application"
LABEL org.opencontainers.image.version="1.0.0"
LABEL org.opencontainers.image.created="2025-04-22"

# Copy SSL certificates from Go image for HTTPS requests
# This is needed to allow the application to run over https, but I unfortunately didn't have time to set this up properly
# COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy executable file from previous stage
COPY --from=builder /app/weather-app /weather-app

# Set working directory
WORKDIR /

# Open port
EXPOSE ${PORT}

# Ideally, i would need to create a HEALTHCHECK, but that would significantly increase the image size
# The whole point is that GO allows using completely bare scratch for running, as it compiles to a binary
# And to create a HEALTHCHECK, you need some curl or something similar that requires a package manager and file system

# Run the application
ENTRYPOINT ["/weather-app"]