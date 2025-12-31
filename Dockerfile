FROM --platform=$TARGETPLATFORM golang:alpine AS builder

ENV CGO_ENABLED=0
ENV GOOS=linux

WORKDIR /build
COPY go.mod .
COPY go.sum .
RUN go mod tidy
COPY . .
RUN CGO_ENABLED=$CGO_ENABLED go build -ldflags "-extldflags='-static'" -o pastecode

FROM alpine:3.23.2
WORKDIR /app
COPY --from=builder /build/pastecode /app/pastecode
COPY --from=builder /build/templates /app/templates
COPY --from=builder /build/static /app/static
ENTRYPOINT ["/app/pastecode"]
