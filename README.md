# Resa

Site web permettant de gérer une liste d'invités à un évenement.

Pour s'inscrire il faut un voucher (code). Le site dispose d'une page d'administration permettant de consulter la liste des invités et gérer les vouchers. Chaque utilisateur possède une page lui permettant de consulter et d'imprimer son invitation.

## Installation

**Prérequis**

* Mingw (64) : [téléchargement](https://sourceforge.net/projects/mingw-w64/files/Toolchains%20targetting%20Win32/Personal%20Builds/mingw-builds/installer/mingw-w64-install.exe/download)
* Go : [téléchargement](https://golang.org/dl/)
* Git : [téléchargement](https://git-scm.com/downloads)


Faire le clone du projet dans :

```
%userprofile%\go\src\
```

```
go get ./...
go install
go run main.go
```

* _go get_ télécharge (clone les dépots) et compile toutes les dépendances qui se trouvent dans les import

* _go install_ Compile et installe toutes les dépendances (rends la compilation plus rapide)

* _go run_ compile et lance le programme

## Création de la base de donnée et 1er lancement

Créer la base de donnée et créer un admin pour pouvoir commencer.
```
go run main.go --init
```
Le programme va demander d'entrer des identifiants pour l'administrateur. Il va ensuite créer la base (des erreurs vont être affiché car des DROP TABLE sont lancés), insérer l'administrateur et un invite par défaut.

On peut ensuite se connecter sur [localhost:8080/admin](http://localhost:8080/admin) et ajouter un voucher à l'admin.

## Configuration

Il y a 2 façon de gérer la configuration :

* En ligne de commande avec des paramètres. Exemple :

```
go run main.go -port 9000 -database exemple.db
```

* A l'aide d'un fichier ini. Exemple (config.ini) :

```ini
port = 9000
database = exemple.db
```

```
go run main.go -config config.ini
```

Liste des paramètres et paramètres par défauts :

```
go run main.go -help
```

## Documentation

```
godoc -http=:6060
```
[http://localhost:6060](http://localhost:6060)