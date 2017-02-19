## Readme Server

Attention, les données contenue dans data doivent être convertis en base64 et mis en tant que string en tant que value pour data

### Format des Requêtes :
#### Fetch la board ####
{"type":"Fetch","dataType":"","data":null}

réponse :
{"type":"Fetch","dataType":"Board","data":"myBase64Data"}

#### Fetch le joueur
{"type":"FetchPlayer","dataType":"","data":null}

réponse :
{"type":"FetchPlayer","dataType":"Player","data":"myBase64Data"}

#### Placer le coup (mettre une origin et changer les paramètres de la pièces)

{"type":"PlacePiece","dataType":"Piece","data":"myBase64Data"}
myBase64Data => {"id":18,"cubes":[{"X":0,"Y":0},{"X":0,"Y":1},{"X":1,"Y":1},{"X":2,"Y":1},{"X":2,"Y":2}],"origin":{"X":10,"Y":10,"playerId":null},"rotation":"N","flipped":false,"playerId":0}

réponse : 

{"type":"Refresh","dataType":"Board","data":"myBase64Data"}


