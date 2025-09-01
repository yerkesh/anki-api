CREATE TABLE users(
    id SERIAL NOT NULL,
    email TEXT NOT NULL,
    username TEXT NOT NULL,
    role TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE collections (
    id SERIAL NOT NULL,
    user_id INT NOT NULL,
    name TEXT NOT NULL,
    native_language TEXT NOT NULL,
    learning_language TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE flashcards (
                            id SERIAL NOT NULL,
                            collection_id INT NOT NULL,
                            frontside TEXT NOT NULL,
                            review_status TEXT NOT NULL,
                            description TEXT NOT NULL,
                            translated TEXT NOT NULL,
                            part_of_speech TEXT NOT NULL,
                            plural_noun TEXT NOT NULL,
                            singular_noun TEXT NOT NULL,
                            base_verb TEXT NOT NULL,
                            past_verb TEXT NOT NULL,
                            participle_verb TEXT NOT NULL,
                            created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                            updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                            repeated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE meanings (
                          id SERIAL NOT NULL,
                          flashcard_id INT NOT NULL,
                          collection_id INT NOT NULL,
                          meaning TEXT NOT NULL,
                          example TEXT NOT NULL
);


CREATE TABLE synonyms (
                         id SERIAL NOT NULL,
                         flashcard_id INT NOT NULL,
                         collection_id INT NOT NULL,
                         word TEXT NOT NULL,
                         translated TEXT NOT NULL
);

CREATE TABLE antonyms (
                         id SERIAL NOT NULL,
                         flashcard_id INT NOT NULL,
                         collection_id INT NOT NULL,
                         word TEXT NOT NULL,
                         translated TEXT NOT NULL
);

ALTER TABLE flashcards
    ADD COLUMN is_deleted BOOLEAN DEFAULT FALSE;