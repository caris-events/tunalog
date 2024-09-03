FROM golang:1.23-alpine AS builder
WORKDIR /workspace

COPY ./ ./
RUN go build -o tunalog .

FROM scratch AS runnable
WORKDIR /var/lib/tunalog

COPY --from=builder /workspace/tunalog /tunalog

EXPOSE 8080
ENTRYPOINT ["/tunalog"]
