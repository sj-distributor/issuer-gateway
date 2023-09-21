## Issuer Gateway

A High Performance Gateway with Certificate Issuance and Flexible Proxy

```
docker run -p 9090:9090  -v $(pwd)/conf/values.yaml:/app/conf/values.yaml da78e64039dd  /app/ig issuer  -f ./conf/values.yaml
```

```
docker run -p 80:80 -p 443:443  -v $(pwd)/conf/values.yaml:/app/conf/values.yaml da78e64039dd  /app/ig gateway  -f ./conf/values.yaml
```