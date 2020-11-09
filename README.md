# auth-go

[![GoDoc](https://godoc.org/github.com/Arveto/auth-go?status.svg)](https://godoc.org/github.com/Arveto/auth-go)

Authentication client application for people or bot.

```txt
id: string
pseudo: string
email: string
level: enum { no, candidate, visitor, standard, administrator }
bot: bool
teams: []string
```

```bash
go get -v -u github.com/Arveto/auth-go
```

# Provider

| URL                          | Comments                                         |
| :--------------------------- | :----------------------------------------------- |
| /auth?app={APP}&r={REDIRECT} | Gen teken and redirect user to `/login?jwt=$JWT` |
| /avatar                      | Get the user's avatar (picture)                  |
| /publickey                   | Get in PEM the provider's public key             |

# Client Application

| URL                           | Comments                                                                                 |
| :---------------------------- | :--------------------------------------------------------------------------------------- |
| /login?r={REDIRECT}           | Redirect the user to the provider                                                        |
| /login?jwt={JWT}&r={REDIRECT} | Save the tocken into a cookie                                                            |
| /avatar?u={ID}                | Redirect to the provider to get the avatar                                               |
| /forget?jwt={JWT}             | (doit être un administrateur) Remove a user (must be an administrator) (`DELETE` method) |
| /me                           | Get connected user information and expiration date (JSON)                                |

The token can be in a cookie or in a HTTP header: `Authorization: Bearer <JWT>`
