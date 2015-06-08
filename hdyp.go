package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/user"
)

var t *template.Template
var staticService = http.FileServer(http.Dir("www"))

func init() {
	t = template.Must(template.ParseGlob("template/*.html"))
	http.HandleFunc("/", handle)
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/logout", handleLogout)
	http.HandleFunc("/_ah/start", serverStart)
}

func main() {
	appengine.Main()
}

func serverStart(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	if appengine.IsDevAppServer() {
		if !GetSettingBool(c, "alreadyInit") {
			if err := loadSomeData(c); err != nil {
				log.Printf(err.Error())
			} else {
				SetSettingBool(c, "alreadyInit", true)
				fmt.Fprint(w, "OK")
				return
			}
		}
	}
	fmt.Fprint(w, "OK")
	return
}

func notFound(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)

	t.ExecuteTemplate(w, "404.html", map[string]interface{}{
		"user": u,
	})
}

func handle(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)

	if r.URL.Path == "/" {
		t.ExecuteTemplate(w, "index.html", map[string]interface{}{
			"user": u,
		})
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

	word, err := GetWord(c, reqWord)
	if err != nil {
		if err == datastore.ErrNoSuchEntity {
			t.ExecuteTemplate(w, "new.html", map[string]interface{}{
				"user": u,
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
			"user":        u,
			"word":        word,
			"noautofocus": true,
		})
		return
	}

	notFound(w, r)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	redir := r.URL.Query().Get("r")
	if redir == "" {
		redir = r.Referer()
		if redir == "" {
			redir = "/"
		}
	}
	c := appengine.NewContext(r)
	u := user.Current(c)
	if u == nil {
		loginURL, err := user.LoginURL(c, redir)
		if err != nil {
			// something bad happened?  FIXME: handle this error some other way
			fmt.Fprint(w, err.Error())
			return
		}
		http.Redirect(w, r, loginURL, http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, redir, http.StatusSeeOther)
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	redir := r.URL.Query().Get("r")
	if redir == "" {
		redir = r.Referer()
		if redir == "" {
			redir = "/"
		}
	}
	c := appengine.NewContext(r)
	u := user.Current(c)
	if u != nil {
		logoutURL, err := user.LogoutURL(c, redir)
		if err != nil {
			// something bad happened?  FIXME: handle this error some other way
			fmt.Fprint(w, err.Error())
			return
		}
		http.Redirect(w, r, logoutURL, http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, redir, http.StatusSeeOther)
}
