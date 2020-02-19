FROM golang:1.13
WORKDIR /tmp/myscheduler-extender
COPY . .
RUN go build -ldflags "-linkmode external -extldflags -static" -a myscheduler-extender.go

FROM scratch
COPY --from=0 /tmp/myscheduler-extender/myscheduler-extender /myscheduler-extender
CMD ["/myscheduler-extender"]
EXPOSE 8888

