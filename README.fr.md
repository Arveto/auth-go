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

# Fournisseur

| URL                          | Commentaire                                                     |
| :--------------------------- | :-------------------------------------------------------------- |
| /auth?app={APP}&r={REDIRECT} | Génère un jeton et renvoie l'utilisateur vers `/login?jwt=$JWT` |
| /avatar                      | Donne l'avatar d'un utilisateur (image)                         |
| /publickey                   | Donne en PEM la clé publique du fournisseur                     |

# Application cliente

| URL                           | Commentaire                                                                               |
| :---------------------------- | :---------------------------------------------------------------------------------------- |
| /login?r={REDIRECT}           | Renvoie l'utilisateur vers le fournisseur d'authentification                              |
| /login?jwt={JWT}&r={REDIRECT} | Enregistre le jeton comme cookie d'authentification                                       |
| /avatar?u={ID}                | Renvoie vers le fournisseur d'authentification pour avoir un avatar                       |
| /forget?jwt={JWT}             | Supprime un utilisateur (administrateur) (méthode `DELETE`)                               |
| /me                           | Donne les informations sur l'utilisateur connecté et la date d'expiration du jeton (JSON) |

Le jeton d'authentification peut être enregistré dans un cookie ou dans un entête HTTP\: `Authorization: Bearer <JWT>`
