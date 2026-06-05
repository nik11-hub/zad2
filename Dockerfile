# syntax=docker/dockerfile:1.7

# --- ETAP 1: BUILDER ---
FROM --platform=$BUILDPLATFORM golang:1.22-alpine AS builder

ARG TARGETOS
ARG TARGETARCH
# Adres repozytorium przekazywany jako build-arg
ARG REPO_URL=git@github.com:nik11-hub/lab-docker-weather-app.git

WORKDIR /src

# Instalacja niezbędnych narzędzi
RUN apk add --no-cache git openssh-client ca-certificates upx

# Konfiguracja SSH dla GitHub
RUN mkdir -p -m 0700 ~/.ssh && ssh-keyscan github.com >> ~/.ssh/known_hosts

# Klonowanie repozytorium z wykorzystaniem montowania klucza SSH
RUN --mount=type=ssh \
    git clone --depth 1 ${REPO_URL} .

# Pobieranie zależności z wykorzystaniem cache
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Kompilacja aplikacji (z wykorzystaniem cache kompilatora)
RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build \
    -ldflags="-s -w" \
    -trimpath \
    -o /app main.go

# Kompresja binarki dla minimalizacji rozmiaru
RUN upx -9 /app

# Przygotowanie struktury certyfikatów dla obrazu scratch
RUN mkdir -p /rootfs/etc/ssl/certs && \
    cp /etc/ssl/certs/ca-certificates.crt /rootfs/etc/ssl/certs/ && \
    cp /app /rootfs/app

# --- ETAP 2: FINAL ---
FROM scratch

# Metadane zgodne z OCI
LABEL org.opencontainers.image.authors="Mikita Liaiko" \
      org.opencontainers.image.description="Minimalistyczna aplikacja pogodowa - Multiarch" \
      org.opencontainers.image.source="https://github.com/nik11-hub/lab-docker-weather-app"

# Kopiowanie przygotowanej zawartości
COPY --from=builder /rootfs /

# Użytkownik nieuprzywilejowany
USER 10000:10000

EXPOSE 8080


HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ["/app", "-health"]

ENTRYPOINT ["/app"]
