FROM node:20-alpine AS bun-builder
WORKDIR /app
ARG UMAMI_ID=""
ARG UMAMI_URL=""
ARG UMAMI_DOMAINS=""
ARG PUBLIC_GA_ID=""

ENV UMAMI_ID=${UMAMI_ID} \
    UMAMI_URL=${UMAMI_URL} \
    UMAMI_DOMAINS=${UMAMI_DOMAINS} \
    PUBLIC_GA_ID=${PUBLIC_GA_ID}

COPY ./web/package.json ./web/package-lock.json* ./
RUN npm ci
COPY ./web .
RUN chmod +x node_modules/.bin/* || true
RUN npm run build

FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=bun-builder /app/dist ./web/dist/
RUN go generate ./ent
RUN CGO_ENABLED=0 go build -ldflags '-s -w' -o tmail cmd/main.go

FROM alpine AS runner
WORKDIR /app
COPY --from=builder /app/tmail .
RUN apk add --no-cache tzdata

ENV HOST=127.0.0.1
ENV PORT=3000
EXPOSE 3000
CMD ["/app/tmail"]
