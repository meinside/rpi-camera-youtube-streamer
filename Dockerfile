# Dockerfile for Golang application

FROM balenalib/raspberrypi3-debian-golang:latest AS builder

# Working directory outside $GOPATH
WORKDIR /src

# Copy go module files and download dependencies
COPY ./go.* ./
RUN go mod download

# Copy source files
COPY ./ ./

# Build source files statically
RUN CGO_ENABLED=0 go build \
		-installsuffix 'static' \
		-o /app \
		.

FROM balenalib/raspberrypi3-debian:latest AS final

# Copy files from temporary image
COPY --from=builder /app /

# Install essential libraries
RUN apt-get update -y && \
		apt-get install -y apt-utils wget git build-essential libx264-dev libraspberrypi-bin

# Build ffmpeg
# (referenced: https://github.com/meinside/rpi-configs/blob/master/bin/install_ffmpeg.sh)
RUN git clone --depth=1 https://github.com/FFmpeg/FFmpeg.git /tmp/ffmpeg && \
	cd /tmp/ffmpeg && \
	./configure --arch=armel --target-os=linux --enable-gpl --enable-nonfree --enable-libx264 && \
	make -j4 && \
	make install && \
	rm -rf /tmp/ffmpeg

# Copy config file
COPY ./config.json /

# Open ports (if needed)
#EXPOSE 8080
#EXPOSE 80
#EXPOSE 443

# Entry point for the built application
ENTRYPOINT ["/app"]
