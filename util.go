package main

import (
	"net/url"

	"golang.org/x/net/context"
)

func mustParseURL(urlString string) string {
	u, _ := url.Parse(urlString)
	return u.String()
}

func loadSomeData(c context.Context) error {
	word := &Word{
		Word: "GIF",
		Pronunciations: []Pronunciation{
			Pronunciation{
				Rating:        65,
				Pronunciation: "jiff",
				IPA:           "/ˈdʒɪf/",
				Description:   "Choosy developers choose GIF.  Like the peanut butter.",
				Sources: []PronunciationSource{
					PronunciationSource{
						Description: "The original creator of the GIF format, Steve Wilhite at Compuserve, documented this pronunciation.  In 2013 at the Webby Award ceremony, he publicly rejected the alternative pronunciation.",
					},
					PronunciationSource{
						URL:         mustParseURL("http://en.wikipedia.org/wiki/GIF#Pronunciation"),
						Description: "Wikipedia",
					},
					PronunciationSource{
						URL: mustParseURL("http://www.cnn.com/2013/05/22/tech/web/pronounce-gif/"),
					},
					PronunciationSource{
						URL: mustParseURL("http://twitpic.com/csdcxf"),
					},
					PronunciationSource{
						URL: mustParseURL("https://twitter.com/Jif/status/337277962837704705"),
					},
					PronunciationSource{
						URL: mustParseURL("http://www.olsenhome.com/gif/"),
					},
					PronunciationSource{
						URL: mustParseURL("https://www.yahoo.com/tech/did-you-just-say-mem-to-ensure-that-you-dont-85736013339.html"),
					},
				},
			},
			Pronunciation{
				Rating:        35,
				Pronunciation: "g'if",
				IPA:           "/ˈɡɪf/",
				Description:   "Like gift without the T.",
				Sources: []PronunciationSource{
					PronunciationSource{
						Description: "Many people believe that because other short G- words use the hard-G sound, this should too.  The English language is strongly based on the argumentum ad populum; because a large number of people prefer this pronunciation, it has been accepted by most dictionaries.",
					},
					PronunciationSource{
						URL: mustParseURL("http://howtoreallypronouncegif.com/"),
					},
				},
			},
		},
	}
	return SetWord(c, word)
}
