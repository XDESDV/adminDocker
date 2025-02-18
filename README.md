# TP : Découverte de Golang à travers l'administration système et Docker

## Objectifs

- Manipuler Docker avec Golang.

- Lister, arrêter et démarrer des conteneurs Docker via une API REST en Go.

- Créer et exécuter un conteneur Docker depuis un programme Go.

## Pré-requis

- Docker installé sur la machine.

- Package Go "github.com/docker/docker/client" pour interagir avec l'API Docker.

## Rendu attendu

Un serveur API REST en Go permettant de gérer les conteneurs Docker.

## Règles de fonctionnements

### Lister les conteneurs Docker en cours d'exécution

**Objectif**

Créer un handler Go qui liste tous les conteneurs en cours d'exécution avec leur nom, ID et statut.
```
conteneur {
  ID string
  Name string
  Status string
}
```

**Consignes**

Utiliser le package "github.com/docker/docker/client" pour interagir avec Docker.

Afficher la liste des conteneurs actifs.

Créer un endpoint :

```GET /dockers ``` pour lister les dockers actifs.

exemple de sortie

```
[
  {
    ID : 9d8a5f1b2a3e, 
    Name : nginx-container,  
    Statut : Up 3 minutes,
  }, 
  {
    ID : 4d8a5e1b2a6d, 
    Name : nginx-container2,  
    Statut : Up 6 minutes,
  },
]
```

### Arrêter et démarrer un conteneur Docker avec Go

**Objectif**

Créer un handler qui permet d’arrêter et de démarrer un conteneur donné via une API REST.

**Consignes**

Créer deux endpoints :

```POST /dockers/stop/:id ``` pour arrêter un conteneur.

```POST /dockers/start/:id ``` pour le redémarrer.

Afficher un message indiquant le succès ou l’échec de l’opération.

### Lancer un conteneur Nginx avec Go

**Objectif**

Créer un handler Go qui lance un conteneur Nginx et affiche ses logs en temps réel.

**Consignes**

- Vérifier si l’image nginx est disponible localement.

- Si ce n'est pas le cas, la télécharger (docker pull nginx).

- Démarrer un conteneur Nginx et afficher son ID.

- Lire et afficher ses logs en temps réel.

**exemple de sortie**
```
Démarrage du conteneur Nginx...  
ID du conteneur : abc123xyz  
Logs du conteneur :  
127.0.0.1 - - [Date] "GET / HTTP/1.1" 200 612 "-" "Mozilla/5.0"
```

### Superviser les ressources utilisées par les conteneurs Docker

**Objectif**

Créer un programme Go qui affiche l’utilisation CPU et RAM des conteneurs en cours d’exécution.

**Consignes**

- Récupérer la liste des conteneurs en cours d'exécution.

- Afficher l’utilisation CPU et RAM de chaque conteneur via l'API Docker.

Créer un endpoint :

```GET /dockers/:id/ressources ``` pour lister les dockers actifs.

**exemple de sortie**
```
Conteneur : nginx-container  
CPU : 2.5%  
Mémoire : 50 Mo / 512 Mo  
```
