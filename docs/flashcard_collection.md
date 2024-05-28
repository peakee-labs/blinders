# Flashcard and Collection Features

## Flashcard
- `flashcard` is managed by the `flashcard repo`.
- Each `flashcard` can belong to one collection, which is referenced by `collectionID`.
- By default, if no `collectionID` is specified, the flashcard will belong to the `default collection`.
- A flashcard created from an event (`translate event`, `explain event`) will have an ID created from the `event ID`. This field is checked to avoid spam creation of flashcards from practice events suggested by the `practice service`.

## Collection
- Currently, `collection` (including metadata and flashcard) is **aggregated from multiple MongoDB collections**.
- By default, **each user will have one `default collection`**. The ID of the `default collection` will be the `userId` to **ensure each user has only one default collection**.
- Each collection will hold a `viewed` field that contains the IDs of flashcards viewed by users, and a `total` field that contains the IDs of flashcards included in the collection.

## Collection Metadata
- `Collection metadata` is located in the `collection-metadata repo`.
- `Collection metadata` is used to hold general information (excluding flashcard values).

## Areas for Improvement:
- We currently have 2 sources of data (`flashcard repo` and `collection metadata repo`), which makes it hard to synchronize and maintain data consistency.
- We could refactor the code to hold `flashcards` and `collection metadata` in a single document of one MongoDB collection (merging `flashcard repo` and `collection metadata repo`).