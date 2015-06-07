package main

import (
	"html/template"
	"net/http"
	"os"
	"strings"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

var t *template.Template
var staticService = http.FileServer(http.Dir("www"))

func init() {
	t = template.Must(template.ParseGlob("template/*.html"))
	http.HandleFunc("/", handle)
}

func main() {
	appengine.Main()
}

func notFound(w http.ResponseWriter, r *http.Request) {
	t.ExecuteTemplate(w, "404.html", nil)
}

func handle(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	if r.URL.Path == "/" {
		t.ExecuteTemplate(w, "index.html", nil)
		return
	}

	if strings.Contains(r.URL.Path, ".") {
		stat, err := os.Stat("www" + r.URL.Path)
		if err != nil || stat.IsDir() || os.IsNotExist(err) {
			// FIXME: probably handle that err not so blindly...
			notFound(w, r)
			return
		}
		staticService.ServeHTTP(w, r)
		return
	}

	if n := strings.Count(r.URL.Path, "/"); n > 1 {
		notFound(w, r)
		return
	}

	reqWord := strings.TrimLeft(r.URL.Path, "/")

	if reqWord != strings.ToLower(reqWord) {
		// it's not all lowercase, so redirect to the canonical path, which is all lowercase
		http.Redirect(w, r, "/"+strings.ToLower(reqWord), http.StatusFound)
		return
	}

	// word := &Word{
	// 	Word: "GIF",
	// 	Pronunciations: []Pronunciation{
	// 		Pronunciation{
	// 			Rating:        65,
	// 			Pronunciation: "jiff",
	// 			IPA:           "/ˈdʒɪf/",
	// 			Description:   "Choosy developers choose GIF.  Like the peanut butter.",
	// 			Sources: []PronunciationSource{
	// 				PronunciationSource{
	// 					Description: "The original creator of the GIF format, Steve Wilhite at Compuserve, documented this pronunciation.  In 2013 at the Webby Award ceremony, he publicly rejected the alternative pronunciation.",
	// 				},
	// 				PronunciationSource{
	// 					URL:         mustParseURL("http://en.wikipedia.org/wiki/GIF#Pronunciation"),
	// 					Description: "Wikipedia",
	// 				},
	// 				PronunciationSource{
	// 					URL: mustParseURL("http://www.cnn.com/2013/05/22/tech/web/pronounce-gif/"),
	// 				},
	// 				PronunciationSource{
	// 					URL: mustParseURL("http://twitpic.com/csdcxf"),
	// 				},
	// 				PronunciationSource{
	// 					URL: mustParseURL("https://twitter.com/Jif/status/337277962837704705"),
	// 				},
	// 				PronunciationSource{
	// 					URL: mustParseURL("http://www.olsenhome.com/gif/"),
	// 				},
	// 				PronunciationSource{
	// 					URL: mustParseURL("https://www.yahoo.com/tech/did-you-just-say-mem-to-ensure-that-you-dont-85736013339.html"),
	// 				},
	// 			},
	// 		},
	// 		Pronunciation{
	// 			Rating:        35,
	// 			Pronunciation: "g'if",
	// 			IPA:           "/ˈɡɪf/",
	// 			Description:   "Like gift without the T.",
	// 			Sources: []PronunciationSource{
	// 				PronunciationSource{
	// 					Description: "Many people believe that because other short G- words use the hard-G sound, this should too.  The English language is strongly based on the argumentum ad populum; because a large number of people prefer this pronunciation, it has been accepted by most dictionaries.",
	// 				},
	// 				PronunciationSource{
	// 					URL: mustParseURL("http://howtoreallypronouncegif.com/"),
	// 				},
	// 			},
	// 		},
	// 	},
	// }
	// if err := SetWord(c, word); err != nil {
	// 	fmt.Fprint(w, err)
	// 	return
	// }

	word, err := GetWord(c, reqWord)
	if err != nil {
		if err == datastore.ErrNoSuchEntity {
			t.ExecuteTemplate(w, "new.html", map[string]interface{}{
				"word": Word{
					Word: reqWord,
				},
			})
			return
		}
		notFound(w, r)
		return
	}

	if word != nil {
		t.ExecuteTemplate(w, "word.html", map[string]interface{}{
			"noautofocus": true,
			"word":        word,
		})
		return
	}

	notFound(w, r)
}
