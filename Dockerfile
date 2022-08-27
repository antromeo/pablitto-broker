FROM docker.io/library/golang:alpine
COPY . /app
WORKDIR /app
RUN go build 
ENTRYPOINT ["./broker-pablitto"]
CMD ["1883"]
EXPOSE 1883
