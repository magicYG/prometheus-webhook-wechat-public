# Webhook adapter for Prometheus & Send Alert To Wechat Group chat

基于[XionZhao/prom-webhook-wechat](https://github.com/XionZhao/prom-webhook-wechat)改进

## Build

Just type and run: `make build`

Generated in the binary file The `./build` Dir

## Usage

```
usage: prom-webhook-wechat [<args>]


   -web.listen-address ":8060"
      Address to listen on for web interface.

 == WECHAT ==

   -wechat.apiurl
      Custom wechat api url

   -wechat.timeout 5s
      Timeout for invoking wechat webhook.
```

## Exmaple

**Do not add to note that there is behind the token of the capacity(The program will get token by corpid and corpsecret)**

#### Start the single webhook and sent to a single group chat
```
./prom-webhook-wechat -wechat.apiurl=api.weixin.qq.com
```
#### Start multiple webhook and sent to multiple group chat
```
./prom-webhook-wechat -config.file=config.yaml -wechat.apiurl=api.weixin.qq.com
```

## Test request prom-webhook-wechat

To view `exmple/send_alert.sh`
