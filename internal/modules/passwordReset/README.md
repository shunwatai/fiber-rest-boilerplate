# password reset

## Workflow

1. Request for the reset password email
   sample request link: `http://<server host>/password-resets/forgot`

API will be called after click "send"
```
POST /password-resets/send
```

request JSON
```
{
  "email": "xxx@exmaple.com"
}
```

2. Open the reset-password link from mailbox
   sample reset link: `http://<server host>/password-resets?token=<token>&userId=<userId>&email=<email>`

API will be called after click "submit"
```
PATCH /password-resets
```

request JSON
```
{
  "password": "new-password"
}
```
