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
	Domain string `json:"domain" validate:"required,hostname_rfc1123"`
	Email string `json:"email" validate:"required,email"`
	Target string `json:"target" validate:"required"`
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

### 4. "上传证书"

1. route definition

- Url: /api/cert/upload
- Method: POST
- Request: `AddCertFormUploadReq`
- Response: `AddOrRenewCertificateResp`

2. request definition



```golang
type AddCertFormUploadReq struct {
	Id uint64 `json:"id"`
	Certificate string `json:"certificate" validate:"required"`
	PrivateKey string `json:"private_key" validate:"required"`
	IssuerCertificate string `json:"issuer_certificate"`
}
```


3. response definition



```golang
type AddOrRenewCertificateResp struct {
}
```

### 5. "重新申请证书"

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

### 6. "证书分页"

1. route definition

- Url: /api/certs
- Method: GET
- Request: `GetCertsPagingReq`
- Response: `GetCertsPagingResp`

2. request definition



```golang
type GetCertsPagingReq struct {
	Page int `form:"page" validate:"required"`
	Size int `form:"size" validate:"required"`
	Domain string `form:"domain,optional"`
	Email string `form:"email,optional"`
}
```


3. response definition



```golang
type GetCertsPagingResp struct {
	Certs []CertDto `json:"certs"`
	Total uint64 `json:"total"`
}
```

