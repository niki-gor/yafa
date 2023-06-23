FROM golang:1.20-alpine AS builder
WORKDIR /app
COPY ./ ./
RUN go mod download
RUN go build -o api cmd/main.go


FROM postgres:15-alpine
COPY --from=builder /app /
COPY ./db/db.sql ./scripts/api.sh /docker-entrypoint-initdb.d/
ENV POSTGRES_USER myuser
ENV POSTGRES_PASSWORD mypassword
ENV POSTGRES_DB forum