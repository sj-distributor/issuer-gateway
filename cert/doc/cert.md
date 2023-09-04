### 1. "用户登录"

1. route definition

- Url: /api/user/login
- Method: POST
- Request: `UserLoginReq`
- Response: `UserLoginResp`

2. request definition



```golang
type UserLoginReq struct {
	Name string `json:"name"  validate:"required,email"`
	Pass string `json:"pass"  validate:"required"`
}
```


3. response definition



```golang
type UserLoginResp struct {
	Token string `json:"token"`
}
```

### 2. "绑定域名"

1. route definition

- Url: /api/domain
- Method: POST
- Request: `AddDomainReq`
- Response: `AddOrRenewCertificateResp`

2. request definition



```golang
type AddDomainReq struct {
	Domain string `json:"domain"  validate:"required"`
	Email string `json:"email"  validate:"required"`
}
```


3. response definition



```golang
type AddOrRenewCertificateResp struct {
}
```

### 3. "申请证书"

1. route definition

- Url: /api/cert
- Method: POST
- Request: `CertificateRequest`
- Response: `AddOrRenewCertificateResp`

2. request definition



```golang
type CertificateRequest struct {
	Id uint64 `json:"id"`
}
```


3. response definition



```golang
type AddOrRenewCertificateResp struct {
}
```

### 4. "重新申请证书"

1. route definition

- Url: /api/cert
- Method: PUT
- Request: `CertificateRequest`
- Response: `AddOrRenewCertificateResp`

2. request definition



```golang
type CertificateRequest struct {
	Id uint64 `json:"id"`
}
```


3. response definition



```golang
type AddOrRenewCertificateResp struct {
}
```

### 5. "证书同步"

1. route definition

- Url: /api/cert/sync
- Method: PUT
- Request: `CertSyncReq`
- Response: `CertSyncResp`

2. request definition



```golang
type CertSyncReq struct {
	Maximum uint64 `form:"maximum"`
}
```


3. response definition



```golang
type CertSyncResp struct {
	Certs []Cert `json:"certs"`
}
```

