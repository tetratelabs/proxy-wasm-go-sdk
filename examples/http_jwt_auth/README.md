## http_jwt_auth

This example demonstrates how to authorize the incoming request depending on jwt token.

In this example, the jwt token should be signed in HMAC256 (HMAC-SHA256). 

```
$ curl -XGET localhost:18000 -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.e30.t-IDcSemACt8x4iTMCda8Yhe3iZaWbvV5XKSTbuAn0M"
example body
```
