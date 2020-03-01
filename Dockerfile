FROM golang:latest AS builder
ADD . /app
WORKDIR /app/company-api
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w" -a -o /main ./src

FROM node:alpine AS node_builder
COPY --from=builder /app/company-web ./
RUN npm install
RUN npm run build

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /main ./
COPY --from=node_builder /dist ./web
RUN chmod +x ./main
EXPOSE ${PORT}
CMD ./main