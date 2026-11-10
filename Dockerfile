FROM docker.akamai.com/images/golang:1.26.9-alpine3.23

# Install the certs
COPY certs /certs

# Set up APK repositories
RUN echo "https://repos.akamai.com/artifactory/upstream-alpine/v3.23/main/" > /etc/apk/repositories && \
    echo "https://repos.akamai.com/artifactory/upstream-alpine/v3.23/community/" >> /etc/apk/repositories
ENV SSL_CLIENT_CERT_FILE=/certs/crt.pem \
    SSL_CLIENT_KEY_FILE=/certs/key.pem
COPY certs/ca.pem /usr/local/share/ca-certificates

# Ensure CA certificates are recognized
# NOTE: Our ca.pem contains multiple certificates; update-ca-certificates will warn and skip the file,
# but OpenSSL still trusts it, so apk works correctly.
RUN apk --no-cache add ca-certificates
RUN update-ca-certificates

# Install necessary packages
RUN apk --no-cache add make bash curl python3 git

# Install go-proxy
COPY proxy/go_proxy.py /usr/local/bin/go-proxy
RUN chmod +x /usr/local/bin/go-proxy

# Install Go tools (versions are defined in the Makefile)
WORKDIR /tmp/tools
COPY Makefile ./
RUN make install-tools BIN=/go/bin GO=go-proxy

# Pre-fetch dependencies to speed up builds
WORKDIR /tmp/go-mod-download
COPY go.mod go.sum ./
RUN go-proxy mod download
