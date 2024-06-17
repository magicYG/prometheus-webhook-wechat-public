FROM alpine:3.9
COPY .build/linux-amd64/prom-webhook-wechat /root/
ENTRYPOINT [ "/root/prom-webhook-wechat" ]