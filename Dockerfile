
FROM node:18-alpine AS frontend-builder
WORKDIR /app/frontend

COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci

COPY frontend/ .
RUN npm run build

FROM golang:1.24-alpine AS backend-builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o socialnet-app main.go


FROM alpine:latest
WORKDIR /app

RUN apk --no-cache add ca-certificates

COPY --from=backend-builder /app/socialnet-app .
COPY --from=frontend-builder /app/frontend/dist ./frontend/dist

ENV FRONTEND_DIR=/app/frontend/dist
ENV UPLOAD_DIR=/tmp/uploads
ENV DB_PATH=/tmp/socialnet.db
ENV SERVER_PORT=8080

EXPOSE 8080

CMD ["./socialnet-app"]