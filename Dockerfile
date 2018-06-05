# Builder environment for the project
FROM golang AS builder

ADD . /go/src/github.com/syntaqx/elkish
WORKDIR /go/src/github.com/syntaqx/elkish

RUN make dep
RUN make build

# Build a scratch container for the binary
FROM scratch
COPY --from=builder /go/src/github.com/syntaqx/elkish/bin/elkish /bin/elkish
RUN touch /var/log/access.log
CMD ["/bin/elkish"]
