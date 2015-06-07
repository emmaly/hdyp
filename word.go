package main

import (
	"strings"

	"code.google.com/p/go-uuid/uuid"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

// Word is the word
type Word struct {
	Word           string
	Pronunciations []Pronunciation `datastore:"-"`
}

// Pronunciation is a pronunciation of the word
type Pronunciation struct {
	UUID          string
	Word          string
	Rating        int
	Pronunciation string
	IPA           string
	Description   string
	Sources       []PronunciationSource
}

// PronunciationSource is a source for the pronunciation, such as URL or otherwise
type PronunciationSource struct {
	URL         string
	Description string
}

// wordKey returns the key for the word
func wordKey(c context.Context, reqWord string) *datastore.Key {
	return datastore.NewKey(c, "Word", strings.ToLower(reqWord), 0, nil)
}

// pronunciationKey returns the key for the pronunciation of the word
func pronunciationKey(c context.Context, reqWord string, pronunciationUUID string) *datastore.Key {
	return datastore.NewKey(c, "Pronunciation", strings.ToUpper(pronunciationUUID), 0, wordKey(c, reqWord))
}

// GetWord fetches a word from the db
func GetWord(c context.Context, reqWord string) (*Word, error) {
	word := new(Word)
	err := datastore.RunInTransaction(c, func(c context.Context) error {
		var err error

		// get the word itself
		key := wordKey(c, reqWord)
		if err = datastore.Get(c, key, word); err != nil {
			return err
		}

		// get the pronunciations
		q := datastore.NewQuery("Pronunciation").Ancestor(key).Order("-Rating").Limit(20) // TODO: paginate beyond limit?
		if _, err := q.GetAll(c, &word.Pronunciations); err != nil {
			return err
		}

		return nil
	}, nil)

	if err != nil {
		return nil, err
	}
	return word, nil
}

// SetWord creates/updates a word in the db
func SetWord(c context.Context, word *Word) error {
	return datastore.RunInTransaction(c, func(c context.Context) error {
		// put the word itself (doesn't include pronunciations)
		if _, err := datastore.Put(c, wordKey(c, word.Word), word); err != nil {
			return err
		}

		// iterate through and put the pronunciations
		for _, pronunciation := range word.Pronunciations {
			if pronunciation.UUID == "" {
				pronunciation.UUID = uuid.New()
			}
			if pronunciation.Word == "" {
				pronunciation.Word = word.Word
			}
			if _, err := datastore.Put(c, pronunciationKey(c, word.Word, pronunciation.UUID), &pronunciation); err != nil {
				return err
			}
		}

		return nil
	}, nil)
}
