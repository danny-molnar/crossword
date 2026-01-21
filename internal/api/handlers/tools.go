package handlers

import (
	"net/http"
	"strconv"
)

func (h *Handler) Anagram(w http.ResponseWriter, r *http.Request) {
	letters := r.URL.Query().Get("letters")
	lenStr := r.URL.Query().Get("len")

	length := 0
	if lenStr != "" {
		n, err := strconv.Atoi(lenStr)
		if err != nil || n < 0 {
			writeErr(w, http.StatusBadRequest, "invalid len")
			return
		}
		length = n
	}

	res, err := h.wl.Anagrams(letters, length)
	if err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, res)
}

func (h *Handler) Pattern(w http.ResponseWriter, r *http.Request) {
	pattern := r.URL.Query().Get("pattern")
	lenStr := r.URL.Query().Get("len")

	length := 0
	if lenStr != "" {
		n, err := strconv.Atoi(lenStr)
		if err != nil || n < 0 {
			writeErr(w, http.StatusBadRequest, "invalid len")
			return
		}
		length = n
	}

	res, err := h.wl.PatternMatch(pattern, length)
	if err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, res)
}
