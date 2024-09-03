FROM golang:1.23-alpine AS builder

WORKDIR /workspace

COPY ./ ./

RUN go build -o tunalog .


FROM scratch AS runnable

# make blog data in a separate folder
WORKDIR /blog

COPY --from=builder /workspace/tunalog /tunalog

EXPOSE 8080
VOLUME [ "/blog" ]
ENTRYPOINT ["/tunalog"]
