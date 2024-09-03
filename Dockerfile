FROM golang:1.23-alpine AS builder

WORKDIR /workspace

COPY ./ ./

RUN go build -o tunalog .


FROM scratch AS runnable

COPY --from=builder /workspace/tunalog /tunalog

ENTRYPOINT ["/tunalog"]

