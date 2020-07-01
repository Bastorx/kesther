# =====================
# Target 'build-env'
# =====================
FROM golang:alpine AS build-env
RUN apk --no-cache add build-base git gcc
ARG GITLAB_TOKEN="${GITLAB_TOKEN}"
RUN git config --global url.https://oauth2:"${GITLAB_TOKEN}"@gitlab.citodi.com/.insteadOf https://gitlab.citodi.com/
COPY . /src
WORKDIR /src
RUN go get ./...
RUN go build -o goapp

# =====================
# Target 'prod'
# =====================
FROM alpine as PROD
WORKDIR /app
ENV PORT 80
COPY --from=build-env /src/goapp /app/
COPY openapi /app/openapi
COPY bin/reset.sh  /app/reset
ENTRYPOINT ./goapp
EXPOSE 80
