

```bash
#docker compose up
pipx run chromadb run --path ./chroma
```


## score threshold

Dans le cadre de la recherche d'information et de la classification automatique :

Le seuil de score est une valeur numérique utilisée pour classifier des documents ou des points de données en fonction de leur pertinence ou de leur similarité avec une requête donnée.

**Fonctionnement**: Un algorithme calcule un score pour chaque document ou point de données, représentant sa pertinence par rapport à la requête. Ces scores sont ensuite comparés au seuil de score prédéfini. Les éléments dont le score est supérieur ou égal au seuil sont considérés comme pertinents et retenus pour une analyse plus approfondie ou une action ultérieure, tandis que ceux dont le score est inférieur au seuil sont écartés.

**Détermination du seuil de score**: Le choix du seuil de score est crucial pour l'efficacité de l'algorithme. Un seuil trop élevé peut entraîner la perte d'informations pertinentes, tandis qu'un seuil trop bas peut conduire à une surcharge d'informations non pertinentes. La définition optimale du seuil dépend de l'application spécifique et des objectifs de la recherche ou de la classification.


Le seuil de score peut être une proportion comprise entre 0 et 1, représentant la proportion de similarité minimale requise entre un document et une requête pour être considéré comme pertinent. Par exemple, un seuil de 0,8 pourrait signifier que seuls les documents ayant une similarité de 80% ou plus avec la requête seront retenus.

Le seuil de score peut également être une valeur numérique définie en fonction d'une mesure de similarité spécifique. Par exemple, si la similarité est mesurée par la distance cosinus, un seuil de 0,5 pourrait indiquer que seuls les documents dont la distance cosinus avec la requête est inférieure ou égale à 0,5 seront considérés comme pertinents.
