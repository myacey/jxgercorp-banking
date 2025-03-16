# build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# install dependencies for caching
COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG SERVICE_NAME
ARG SERVICE_PATH="./services/${SERVICE_NAME}"

RUN echo ${SERVICE_NAME}
RUN echo ${SERVICE_PATH}

# build binaries
RUN go build -o /app/bin/${SERVICE_NAME} ${SERVICE_PATH}/cmd/main.go

# final stage
FROM alpine:latest

ARG SERVICE_NAME
ENV DEPLOY_COMMAND="./${SERVICE_NAME}"

COPY --from=builder /app/bin/${SERVICE_NAME} .
COPY .env .

ARG PORT
EXPOSE ${PORT}

RUN chmod +x ${SERVICE_NAME}

CMD ${DEPLOY_COMMAND}