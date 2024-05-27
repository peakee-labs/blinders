# Flashcard and collection feature


# flashcard
- `flashcard` is managed by `flashcard repo`
- each `flashcard` could be belongs to 1 collection which is referenced at `collectionID`
- by default flashcard when created if not specifies any collection will belong to `default collection`
- flashcard that created from event (`translate event`, `explain event`) will have ID that created from `event ID`, this field will be check to avoid spam create flashcard from practice event that is suggested by `practice service`

# collection
- currently, `collection` (include metadata, flashcard) is **aggregated from multiple mongo collection**
- by default, **user will have 1 `default collection`**, ID of `default collection` will be `userId` to **make sure **each user will have only 1 default collection**
- each collection will hold *`viewed` fields that contains ID of flashcards that are viewed by users*, and *`total` fields that contains ID of flashcards included in the collection*.

# collection metadata
- `collection metadata` located in `collection-metadata repo`
- `collection metadata` is used to hold general information (without flashcards value)


# Improve area:
- We currently have 2 source of data (`flashcard repo` and `collection metadata repo`), that makes it hard to sync and make data consistent.
- We could refactor the code, hold `flashcards` and `collection metadata` in single document of 1 mongo collection (merge `flashcard repo` and `collection metadata repo`)