package server

import "net/http"

func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/next", http.StatusSeeOther)
}
