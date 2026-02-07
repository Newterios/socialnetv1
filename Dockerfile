
FROM node:18-alpine AS frontend-builder
WORKDIR /app/frontend

COPY frontend/package.json frontend/package-lock.json ./
RUN npm install

COPY frontend/ .
RUN npm run build

FROM golang:1.23-alpine AS backend-builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o socialnet-app main.go


FROM alpine:latest
WORKDIR /root/

RUN apk --no-cache add ca-certificates

COPY --from=backend-builder /app/socialnet-app .

COPY --from=frontend-builder /app/frontend/dist ./frontend/dist


EXPOSE 8080


CMD ["./socialnet-app"]