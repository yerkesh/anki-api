package handlers

type OptionsSetter func(o *Options)

func NewOptions(setters ...OptionsSetter) Options {
	o := Options{}
	for _, setter := range setters {
		setter(&o)
	}

	return o
}

func WithValidator(usecase structValidator) OptionsSetter {
	return func(o *Options) {
		o.validator = usecase
	}
}

func WithFlashcards(usecase flashcardUsecase) OptionsSetter {
	return func(o *Options) {
		o.flashcard = usecase
	}
}

func WithCollections(usecase collectionUsecase) OptionsSetter {
	return func(o *Options) {
		o.collection = usecase
	}
}

func WithUsers(usecase userUsecase) OptionsSetter {
	return func(o *Options) {
		o.user = usecase
	}
}
