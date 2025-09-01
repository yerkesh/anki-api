package entity

import "time"

// Pagination query parameters
type PageableQueryParams struct {
	Page uint64 `query:"page" validate:"omitempty,min=1"`
	Size uint64 `query:"page_size" validate:"omitempty,min=1"`
}

// Normalize sets default values for the pageable query parameters, if not provided.
func (pqp *PageableQueryParams) Normalize() {
	if pqp.Size == 0 {
		pqp.Size = 10
	}

	if pqp.Page == 0 {
		pqp.Page = 1
	}
}

func (pqp *PageableQueryParams) Offset() uint64 {
	return (pqp.Page - 1) * pqp.Size
}

type PaginationMeta struct {
	TotalItems  uint64 `json:"total_items"`
	TotalPages  uint64 `json:"total_pages"`
	CurrentPage uint64 `json:"current_page"`
	PageSize    uint64 `json:"page_size"`
}

func NewPaginationMeta(totalItems, pageSize, page uint64) PaginationMeta {
	totalPages := (totalItems + pageSize - 1) / pageSize // Round up for total pages

	return PaginationMeta{
		TotalItems:  totalItems,
		TotalPages:  totalPages,
		PageSize:    pageSize,
		CurrentPage: page,
	}
}

type Flashcard struct {
	ID           int32        `json:"id"`
	CollectionID int32        `json:"collection_id" validate:"required"`
	Frontside    string       `json:"frontside" validate:"required"`
	Translated   string       `json:"translated"`
	Description  string       `json:"description"`
	Meanings     []Meaning    `json:"meanings"`
	Synonyms     []Synonym    `json:"synonyms"`
	Antonyms     []Antonym    `json:"antonyms"`
	Grammar      Grammar      `json:"grammar"`
	ReviewStatus ReviewStatus `json:"review_status"`
	RepeatedAt   time.Time    `json:"repeated_at"`
}

type GetFlashcardsResponse struct {
	Flashcards     []Flashcard    `json:"items"`
	PaginationMeta PaginationMeta `json:"meta"`
}

type Synonym struct {
	Word       string `json:"word"`
	Translated string `json:"translated"`
}

type Antonym struct {
	Word       string `json:"word"`
	Translated string `json:"translated"`
}

type Meaning struct {
	Meaning string `json:"meaning"`
	Example string `json:"example"`
}

type Grammar struct {
	PartOfSpeech string    `json:"part_of_speech"`
	NounForms    NounForms `json:"noun_forms"`
	VerbForms    VerbForms `json:"verb_forms"`
}

type NounForms struct {
	Singular string `json:"singular"`
	Plural   string `json:"plural"`
}

type VerbForms struct {
	Base       string `json:"base"`
	Past       string `json:"past"`
	Participle string `json:"participle"`
}

type ReviewStatus string

const (
	HardReviewStatus   ReviewStatus = "hard"
	EasyReviewStatus   ReviewStatus = "easy"
	RepeatReviewStatus ReviewStatus = "repeat"
)

func (r ReviewStatus) String() string {
	return string(r)
}
