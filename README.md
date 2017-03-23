# Readme Server
## Format des Requêtes :
Attention, les données contenue dans data doivent être convertis en base64 et mis en tant que string en tant que value pour data

Si le client envois un callbackId lors d'une requête, il sera aussi présent dans la réponse. Les messages Broadcasté par le serveur comme les messages de type Refresh n'ont pas de callbackId
### Fetch la board
#### requête :
{"type":"Fetch","dataType":"","data":null, "callbackId":""}
#### réponse :
{"type":"Fetch","dataType":"Board","data":"myBase64Data", "callbackId":""}
### Fetch le joueur
#### requête :
{"type":"FetchPlayer","dataType":"","data":null, "callbackId":""}
#### réponse :
{"type":"FetchPlayer","dataType":"Player","data":"myBase64Data", "callbackId":""}
### Placer un coup (mettre une origin et changer les paramètres de la pièces)
#### requête :
{"type":"PlacePiece","dataType":"Piece","data":"myBase64Data", "callbackId":""}

myBase64Data => {"id":18,"cubes":[{"X":0,"Y":0},{"X":0,"Y":1},{"X":1,"Y":1},{"X":2,"Y":1},{"X":2,"Y":2}],"origin":{"X":10,"Y":10,"playerId":null},"rotation":"N","flipped":false,"playerId":0}
#### réponse au joueur (une des deux) :
{"type":"PlacementConfirmed","dataType":"","data":null, "callbackId":""}

{"type":"PlacementRefused","dataType":"","data":null, "callbackId":""}

#### placer un coup automatiquement (pas de réponse envoyé par le serveur) :
{"type":"PlaceRandom","dataType":"","data":null, "callbackId":""}
### Liste des messages Broadcastés à tous les joueurs
#### broadcast à tous les joueurs de la board si un est validé :
{"type":"Refresh","dataType":"Board","data":"myBase64Data", "callbackId":""}
#### Game Over
{"type":"GameOver","dataType":"","data":null, "callbackId":""}

#### Abandonner puis quitter uner partie
{"type":"Concede","dataType":"","data":null, "callbackId":""}

{"type":"Quit","dataType":"","data":null, "callbackId":""}


### Depuis l'état start
{"type":"CreateDemo","dataType":"","data":null, "callbackId":""}

{"type":"Authenticate","dataType":"string","data":"myBase64token", "callbackId":""}

myBase64token => "MyToken"

### Depuis l'état home
{"type":"CreateLobby","dataType":"","data":null, "callbackId":""}

{"type":"JoinLobby","dataType":"int","data":"myInt", "callbackId":""}

{"type":"Broadcast","dataType":"ListLobby","data":"myBase64Data", "callbackId":""}

{"type":"Broadcast","dataType":"ListGame","data":"myBase64Data", "callbackId":""}

### Depuis l'état lobby
{"type":"Start","dataType":"","data":null, "callbackId":""}  => Attention seul le Master peut démarre la partie

{"type":"FetchLobby","dataType":"Lobby","data":null, "callbackId":""}  => Attention le message est implémenté mais le data est vide

{"type":"Sit","dataType":"int","data":"myInt", "callbackId":""}

{"type":"Unsit","dataType":"int","data":"myInt", "callbackId":""}

{"type":"SitAI","dataType":"int","data":"myInt", "callbackId":""} => Attention seul le Master peut utiliser cette commande

{"type":"UnsitAI","dataType":"int","data":"myInt", "callbackId":""} => Attention seul le Master peut utiliser cette commande

{"type":"Quit","dataType":"","data":null, "callbackId":""}

myInt => "myInt" | => 0 : "MA==", 1 : "MQ==", 2 : "Mg==", 3 : "Mw=="

### Depuis l'état game

Si vous êtes en mode demo vous serez renvoyé à l'état start, si vous êtes authentifié et avez démarré la game depuis un lobby vous serez renvoyé dans l'état home.
