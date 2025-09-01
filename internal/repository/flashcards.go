package repository

import (
	"context"
	"fmt"

	"anki-api/internal/entity"
	"anki-api/internal/repository/generated"
)

type FlashcardsRepo struct {
	generated.Querier
}

func NewFlashcardsRepo(querier generated.Querier) *FlashcardsRepo {
	return &FlashcardsRepo{querier}
}

func (f *FlashcardsRepo) CreateFlashcard(ctx context.Context, card entity.Flashcard) (int32, error) {
	// TODO Transaction needed
	id, err := f.InsertFlashcard(ctx, generated.InsertFlashcardParams{
		CollectionID:   card.CollectionID,
		Frontside:      card.Frontside,
		Description:    card.Description,
		Translated:     card.Translated,
		ReviewStatus:   entity.RepeatReviewStatus.String(),
		PartOfSpeech:   card.Grammar.PartOfSpeech,
		PluralNoun:     card.Grammar.NounForms.Plural,
		SingularNoun:   card.Grammar.NounForms.Singular,
		BaseVerb:       card.Grammar.VerbForms.Base,
		PastVerb:       card.Grammar.VerbForms.Past,
		ParticipleVerb: card.Grammar.VerbForms.Participle,
	})
	if err != nil {
		return 0, fmt.Errorf("createFlashcards: couldn't create flashcards, Cause: %w", err)
	}

	argMeanings := generated.InsertMeaningsParams{}

	for _, meaning := range card.Meanings {
		argMeanings.Meanings = append(argMeanings.Meanings, meaning.Meaning)
		argMeanings.Examples = append(argMeanings.Examples, meaning.Example)
		argMeanings.FlashcardIds = append(argMeanings.FlashcardIds, id)
		argMeanings.CollectionIds = append(argMeanings.CollectionIds, card.CollectionID)
	}

	if err = f.InsertMeanings(ctx, argMeanings); err != nil {
		return 0, fmt.Errorf("f.InsertMeanings: couldn't insert meanings, Cause: %w", err)
	}

	argSynonyms := generated.InsertSynonymsParams{}

	for _, synonym := range card.Synonyms {
		argSynonyms.FlashcardIds = append(argSynonyms.FlashcardIds, id)
		argSynonyms.CollectionIds = append(argSynonyms.CollectionIds, card.CollectionID)
		argSynonyms.Word = append(argSynonyms.Word, synonym.Word)
		argSynonyms.Translated = append(argSynonyms.Translated, synonym.Translated)
	}

	if err = f.InsertSynonyms(ctx, argSynonyms); err != nil {
		return 0, fmt.Errorf("f.InsertSynonyms: couldn't insert synonyms, Cause: %w", err)
	}

	argAntonyms := generated.InsertAntonymsParams{}

	for _, antonym := range card.Antonyms {
		argAntonyms.FlashcardIds = append(argAntonyms.FlashcardIds, id)
		argAntonyms.CollectionIds = append(argAntonyms.CollectionIds, card.CollectionID)
		argAntonyms.Word = append(argAntonyms.Word, antonym.Word)
		argAntonyms.Translated = append(argAntonyms.Translated, antonym.Translated)
	}

	if err = f.InsertAntonyms(ctx, argAntonyms); err != nil {
		return 0, fmt.Errorf("f.InsertAntonyms: couldn't insert antonyms, Cause: %w", err)
	}

	return id, nil
}

func (f *FlashcardsRepo) UpdateReview(ctx context.Context, flashcardID int, status entity.ReviewStatus) error {
	if err := f.UpdateFlashcardStatus(ctx, generated.UpdateFlashcardStatusParams{
		ID:           int32(flashcardID),
		ReviewStatus: status.String()},
	); err != nil {
		return fmt.Errorf("f.UpdateReview: %w", err)
	}

	return nil
}

func (f *FlashcardsRepo) GetFlashcard(ctx context.Context, flashcardID int) (entity.Flashcard, error) {
	card, err := f.SelectFlashcard(ctx, int32(flashcardID))
	if err != nil {
		return entity.Flashcard{}, fmt.Errorf("f.SelectFlashcard: %w", err)
	}

	meanings, err := f.SelectMeaningsByCard(ctx, card.ID)
	if err != nil {
		return entity.Flashcard{}, fmt.Errorf("f.SelectMeanings: couldn't get meanings, card.ID: %d Cause: %w", card.ID, err)
	}

	synonyms, err := f.SelectSynonymsByCard(ctx, card.ID)
	if err != nil {
		return entity.Flashcard{}, fmt.Errorf("f.SelectSynonyms: couldn't get synonyms, card.ID: %d Cause: %w", card.ID, err)
	}

	antonyms, err := f.SelectAntonymsByCard(ctx, card.ID)
	if err != nil {
		return entity.Flashcard{}, fmt.Errorf("f.SelectAntonyms: couldn't get antonyms, card.ID: %d Cause: %w", card.ID, err)
	}

	return cardToEntity(card, synonyms, antonyms, meanings), nil
}

func cardToEntity(card generated.SelectFlashcardRow,
	synonyms []generated.SelectSynonymsByCardRow,
	antonyms []generated.SelectAntonymsByCardRow,
	meanings []generated.SelectMeaningsByCardRow,
) entity.Flashcard {
	synonymsEnt := make([]entity.Synonym, 0, len(synonyms))
	antonymsEnt := make([]entity.Antonym, 0, len(antonyms))
	meaningsEnt := make([]entity.Meaning, 0, len(meanings))

	for _, synonym := range synonyms {
		synonymsEnt = append(synonymsEnt, entity.Synonym{
			Word:       synonym.Word,
			Translated: synonym.Translated,
		})
	}

	for _, antonym := range antonyms {
		antonymsEnt = append(antonymsEnt, entity.Antonym{
			Word:       antonym.Word,
			Translated: antonym.Translated,
		})
	}

	for _, meaning := range meanings {
		meaningsEnt = append(meaningsEnt, entity.Meaning{
			Meaning: meaning.Meaning,
			Example: meaning.Example,
		})
	}

	res := entity.Flashcard{
		ID:           card.ID,
		CollectionID: card.CollectionID,
		Frontside:    card.Frontside,
		Description:  card.Description,
		Translated:   card.Translated,
		ReviewStatus: entity.ReviewStatus(card.ReviewStatus),
		Grammar: entity.Grammar{
			PartOfSpeech: card.PartOfSpeech,
			NounForms: entity.NounForms{
				Plural:   card.PluralNoun,
				Singular: card.SingularNoun,
			},
			VerbForms: entity.VerbForms{
				Base:       card.BaseVerb,
				Past:       card.PastVerb,
				Participle: card.ParticipleVerb,
			},
		},
		Meanings: meaningsEnt,
		Synonyms: synonymsEnt,
		Antonyms: antonymsEnt,
	}

	return res
}

func (f *FlashcardsRepo) GetFlashcards(ctx context.Context, collectionID int32, params entity.PageableQueryParams) ([]entity.Flashcard, error) {
	cards, err := f.SelectFlashcards(ctx, generated.SelectFlashcardsParams{
		CollectionID: collectionID,
		Limit:        int32(params.Size),
		Offset:       int32(params.Offset())})
	if err != nil {
		return nil, fmt.Errorf("f.SelectFlashcards: couldn't get flashcards, collectionID: %d Cause: %w", collectionID, err)
	}

	meanings, err := f.SelectMeanings(ctx, collectionID)
	if err != nil {
		return nil, fmt.Errorf("f.SelectMeanings: couldn't get meanings, collectionID: %d Cause: %w", collectionID, err)
	}

	synonyms, err := f.SelectSynonyms(ctx, collectionID)
	if err != nil {
		return nil, fmt.Errorf("f.SelectSynonyms: couldn't get synonyms, collectionID: %d Cause: %w", collectionID, err)
	}

	antonyms, err := f.SelectAntonyms(ctx, collectionID)
	if err != nil {
		return nil, fmt.Errorf("f.SelectAntonyms: couldn't get antonyms, collectionID: %d Cause: %w", collectionID, err)
	}

	return cardsToEntity(cards, synonyms, antonyms, meanings), nil
}

func (f *FlashcardsRepo) GetFlashcardsTotal(ctx context.Context, collectionID int32) (int64, error) {
	total, err := f.SelectFlashcardsTotal(ctx, collectionID)
	if err != nil {
		return 0, fmt.Errorf("f.GetFlashcardsTotal: %w", err)
	}

	return total, nil
}

func cardsToEntity(
	cards []generated.SelectFlashcardsRow,
	synonyms []generated.SelectSynonymsRow,
	antonyms []generated.SelectAntonymsRow,
	meanings []generated.SelectMeaningsRow,
) []entity.Flashcard {
	synonymsMap := make(map[int32][]entity.Synonym)
	antonymsMap := make(map[int32][]entity.Antonym)
	meaningsMap := make(map[int32][]entity.Meaning)

	for _, s := range synonyms {
		synonymsMap[s.FlashcardID] = append(synonymsMap[s.FlashcardID], entity.Synonym{
			Word:       s.Word,
			Translated: s.Translated,
		})
	}

	for _, a := range antonyms {
		antonymsMap[a.FlashcardID] = append(antonymsMap[a.FlashcardID], entity.Antonym{
			Word:       a.Word,
			Translated: a.Translated,
		})
	}

	for _, m := range meanings {
		meaningsMap[m.FlashcardID] = append(meaningsMap[m.FlashcardID], entity.Meaning{
			Meaning: m.Meaning,
			Example: m.Example,
		})
	}

	resCards := make([]entity.Flashcard, 0, len(cards))
	for _, card := range cards {
		resCards = append(resCards, entity.Flashcard{
			ID:           card.ID,
			CollectionID: card.CollectionID,
			Frontside:    card.Frontside,
			Description:  card.Description,
			Translated:   card.Translated,
			ReviewStatus: entity.ReviewStatus(card.ReviewStatus),
			RepeatedAt:   card.RepeatedAt.Time,
			Grammar: entity.Grammar{
				PartOfSpeech: card.PartOfSpeech,
				NounForms: entity.NounForms{
					Plural:   card.PluralNoun,
					Singular: card.SingularNoun,
				},
				VerbForms: entity.VerbForms{
					Base:       card.BaseVerb,
					Past:       card.PastVerb,
					Participle: card.ParticipleVerb,
				},
			},
			Meanings: meaningsMap[card.ID],
			Synonyms: synonymsMap[card.ID],
			Antonyms: antonymsMap[card.ID],
		})
	}

	return resCards
}

func (f *FlashcardsRepo) DeleteFlashcard(ctx context.Context, flashcardID int) error {
	if err := f.DeleteFlashcardSoft(ctx, int32(flashcardID)); err != nil {
		return fmt.Errorf("f.UpdateReview: %w", err)
	}

	return nil
}
