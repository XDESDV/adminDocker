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

Créer un handler Go qui liste tous les conteneurs en cours d'exécution avec les informations de types.Container de ```"github.com/docker/docker/api/types"```.
```
conteneur {
      Id string
      Names []string
      Image string
      ImageID string
      Command string
      Created int
      ...
}
```

**Consignes**

Utiliser le package "github.com/docker/docker/client" pour interagir avec Docker.

Utiliser les réponses API standardisées

**Modeles**
```
// WSResponse is the standardized response format.
// - Meta : *Pre-formatted response header returning data.
// - Data : *Data or list of data returned.
type WSResponse struct {
	Meta MetaResponse `json:"meta"`
	Data interface{}  `json:"data"`
}

// MetaResponse is a valid response header
// - ObjectName : *Information returned to the front end to let it know what format it is receiving.*
// - TotalCount : *Total number of records the request can return.
// - Offset : *Starting position of the list of records returned to the Front.
// - Count : *Number of records returned to the Front.
type MetaResponse struct {
	ObjectName string `json:"object_name"`
	TotalCount int    `json:"total_count"`
	Offset     int    `json:"offSet"`
	Count      int    `json:"count"`
}
```

Afficher la liste des conteneurs actifs.

Créer un endpoint :

```GET /dockers ``` pour lister les dockers actifs.

**exemple de sortie**
```
meta: {
  ObjectName: "dockers",
  TotalCount: 2,
  Offset: 0
  Count: 2
},
data: {
  [
    {
      "Id": "123456789abc",
      "Names": [
        "/fake-nginx"
      ],
      "Image": "",
      "ImageID": "",
      "Command": "",
      "Created": 0,
      "Ports": null,
      "Labels": null,
      "State": "running",
      "Status": "Up 10 minutes",
      "HostConfig": {
        
      },
      "NetworkSettings": null,
      "Mounts": null
    },
    {
      "Id": "987654321xyz",
      "Names": [
        "/fake-redis"
      ],
      "Image": "",
      "ImageID": "",
      "Command": "",
      "Created": 0,
      "Ports": null,
      "Labels": null,
      "State": "exited",
      "Status": "Exited (0) 5 minutes ago",
      "HostConfig": {
        
      },
      "NetworkSettings": null,
      "Mounts": null
    }
  ]
}
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
