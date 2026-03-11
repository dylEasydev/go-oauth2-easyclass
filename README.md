# Serveur oauth2 easy class

## Introduction 

Dans le but de monter une architecture micro service , pour améliorer éducation en ligne et créer un système naturel à l'ensignement collaboratif. Nous avons mis en place un serveur oauth2 pour facilter integration de different clients à ce système on reste dans un monde collaboratif on s'apprend mutuelement. Le code est ouvert et est sous Licences MIT. 

Pour en savoir plus oauth2.0 consulter ce site [apprendre sur oauth2](https://oauth.net) ou demander à un LLM chacun ses préférences.

## Prérequis

Ce qu'il faut pour démarer ce projet , le tester et voir le cloner pour l'améliorer 

>[!IMPORTANT]
> pour ce qui utilise docker une image sera bientôt mis en ligne

- [git](https://git-scm.com)
- [go](https://go.dev)
- [openssl](https://www.openssl.org)

## Mis en oeuvre du projet 

### clonage du dépôt 
Tous d'aborb il faut cloner le dépôt git , ouvrez votre terminal et copier le commande si dessous 

```bash
$ git clone https://github.com/dylEasydev/go-oauth2-easyclass.git
```
### installer toutes les dépendences néccesaires avec go 
```bash
$ go mod tidy
```

### Génération des clés
Allez à la racine de vôtre projet 
```bash
$ cd ./go-oauth2-easyclass
```
en suite créer le dossier pour les clé
```bash
$ mkdir key
```
passons à la génération du certificat autosigné du serveur htps
```bash
$ cd ./key
$ openssl genrsa -out server.key 2048
$ openssl req -new -x509 -days 365 -key server.key -out server.perm
```
la clé RSA du client test 

```bash
$ cd ./key
$ openssl genrsa -out private.key 2048
$ openssl rsq -pubout -in private.key -out public.key
```

### Remplicage du fichier d'environement

Toujours à la racine du projet créer un fichier `.env` que vous remplissez comme suit 

```env
DB_HOST=
DB_USER=
DB_PASSWORD=
DB_NAME= 
DB_PORT=5432
PORT=3000
USER_NAME=
COMPANING_MAIl=

//mot de passe solide au moins 8 carractère donc au moins ( lettre miniscule , majuscule , symbole et carractère numériques )
USER_PASSWORD=

//salt pour haché les code de verification
KEY_HASH=

//mail password pour les comptes gmail ici (https://myaccount.google.com/apppasswords)
PASSWORD_MAIL=

//secret du provider ory/fosite (min 32 octects ) 
SECRET=

//secret du client de test 
SECRET_CLIENT=

//secret de rotation du client 
SECRET_CLIENT2=
```

## Démarage du serveur
Depuis la racine du projet lancer 
```bash
$ go run .
```

## test Initial 

```bash
$ curl -X GET \
https://localhost:3000
```

## Documentation
[documentation](./doc/doc.md)
>[!IMPORTANT]
> en cours de rédaction 
