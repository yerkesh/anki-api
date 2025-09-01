-- name: InsertFlashcard :one
INSERT INTO flashcards (collection_id, frontside, review_status, description, translated, part_of_speech, plural_noun, singular_noun, base_verb, past_verb, participle_verb)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id;

-- name: InsertMeanings :exec
INSERT INTO meanings (collection_id, flashcard_id, meaning, example)
SELECT unnest(@collection_ids::integer[]),
       unnest(@flashcard_ids::integer[]),
       unnest(@meanings::text[]),
       unnest(@examples::text[]);

-- name: InsertSynonyms :exec
INSERT INTO synonyms (collection_id, flashcard_id, word, translated)
SELECT unnest(@collection_ids::integer[]),
       unnest(@flashcard_ids::integer[]),
       unnest(@word::text[]),
       unnest(@translated::text[]);

-- name: InsertAntonyms :exec
INSERT INTO antonyms (collection_id, flashcard_id, word, translated)
SELECT unnest(@collection_ids::integer[]),
       unnest(@flashcard_ids::integer[]),
       unnest(@word::text[]),
       unnest(@translated::text[]);

-- name: SelectFlashcards :many
SELECT
    f.id,
    f.collection_id,
    f.frontside,
    f.review_status,
    f.description,
    f.translated,
    f.part_of_speech,
    f.plural_noun,
    f.singular_noun,
    f.base_verb,
    f.past_verb,
    f.participle_verb,
    f.repeated_at
FROM flashcards f
WHERE collection_id = $1 AND is_deleted = FALSE
ORDER BY
    CASE f.review_status
        WHEN 'repeat' THEN 1
        WHEN 'hard'   THEN 2
        WHEN 'easy'   THEN 3
        ELSE 4
        END,
    f.repeated_at
    LIMIT $2
OFFSET $3;


-- name: SelectFlashcard :one
SELECT
    f.id,
    f.collection_id,
    f.frontside,
    f.review_status,
    f.description,
    f.translated,
    f.part_of_speech,
    f.plural_noun,
    f.singular_noun,
    f.base_verb,
    f.past_verb,
    f.participle_verb,
    f.repeated_at
FROM flashcards f
WHERE f.id = $1;

-- name: SelectFlashcardsTotal :one
SELECT COUNT(*) FROM flashcards WHERE collection_id = $1 AND is_deleted = FALSE;

-- name: SelectMeanings :many
SELECT flashcard_id, meaning, example FROM meanings
WHERE collection_id = $1;

-- name: SelectSynonyms :many
SELECT flashcard_id, word, translated FROM synonyms
WHERE collection_id = $1;

-- name: SelectAntonyms :many
SELECT flashcard_id, word, translated FROM antonyms
WHERE collection_id = $1;

-- name: UpdateFlashcardStatus :exec
UPDATE flashcards SET review_status = $1, repeated_at = now() WHERE id = $2;

-- name: DeleteFlashcardSoft :exec
UPDATE flashcards SET is_deleted = TRUE WHERE id = $1;

-- name: SelectMeaningsByCard :many
SELECT flashcard_id, meaning, example FROM meanings
WHERE flashcard_id = $1;

-- name: SelectSynonymsByCard :many
SELECT flashcard_id, word, translated FROM synonyms
WHERE flashcard_id = $1;

-- name: SelectAntonymsByCard :many
SELECT flashcard_id, word, translated FROM antonyms
WHERE flashcard_id = $1;