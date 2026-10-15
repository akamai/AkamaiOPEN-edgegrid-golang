FROM docker.akamai.com/images/golang:1.26.9-alpine3.23

# Install the certs
COPY certs /certs

# Set up APK repositories
RUN echo "https://repos.akamai.com/artifactory/upstream-alpine/v3.22/main/" > /etc/apk/repositories && \
    echo "https://repos.akamai.com/artifactory/upstream-alpine/v3.22/community/" >> /etc/apk/repositories
ENV SSL_CLIENT_CERT_FILE=/certs/crt.pem \
    SSL_CLIENT_KEY_FILE=/certs/key.pem
COPY certs/ca.pem /usr/local/share/ca-certificates

# Ensure CA certificates are recognized
# NOTE: Our ca.pem contains multiple certs; update-ca-certificates will warn and skip,
# but OpenSSL still trusts it, so apk works correctly.
RUN apk --no-cache add ca-certificates
RUN update-ca-certificates

# Install necessary packages
RUN apk --no-cache add make bash curl python3

# Install go-proxy
COPY proxy/go_proxy.py /usr/local/bin/go-proxy
RUN chmod +x /usr/local/bin/go-proxy

# Set GOTOLCHAIN to local for repeatable builds
# Disable Go checksum database as there is no connectivity
ENV GOTOLCHAIN=local \
    GOSUMDB=off

# Install go tools
ARG GOIMPORTS_VERSION=v0.50.0
ARG GOLANGCI_LINT_VERSION=v2.13.2
#ARG GOCOV_VERSION=v1.1.0
#ARG GOCOV_XML_VERSION=v1.1.0
ARG GOCOVER_COBERTURA_VERSION=v1.5.0
ARG GO_JUNIT_REPORT_VERSION=v2.1.0

RUN echo "Installing goimports" && \
    go-proxy install golang.org/x/tools/cmd/goimports@${GOIMPORTS_VERSION} && \
    echo "Installing golangci-lint" && \
    go-proxy install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@${GOLANGCI_LINT_VERSION} && \
    echo "Installing gocov" && \
    go-proxy install github.com/boumenot/gocover-cobertura@${GOCOVER_COBERTURA_VERSION} && \

#    go-proxy install github.com/axw/gocov/gocov@${GOCOV_VERSION} && \
#    echo "Installing gocov-xml" && \
#    go-proxy install github.com/AlekSi/gocov-xml@${GOCOV_XML_VERSION} && \
    echo "Installing go-junit-report" && \
    go-proxy install github.com/jstemmer/go-junit-report/v2@${GO_JUNIT_REPORT_VERSION}

# Pre-fetch dependencies to speed up builds
WORKDIR /tmp/go-mod-download
COPY go.mod go.sum ./
RUN go-proxy mod download
