package main

import (
	"backend/internal/models"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
)

func (app *application) Home(w http.ResponseWriter, r *http.Request) {
	var payload = struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Version string `json:"version"`
	}{
		Status:  "active",
		Message: "Webforum is up and running",
		Version: "1.0.0",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

func (app *application) authenticate(w http.ResponseWriter, r *http.Request) {
	// read json payload
	var requestPayload struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// validate user against database
	user, err := app.DB.GetUserByUsername(requestPayload.Username)
	if err != nil {
		app.errorJSON(w, errors.New("invalid credentials: No username matched!"), http.StatusBadRequest)
		return
	}

	// check password
	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		app.errorJSON(w, errors.New("invalid credentials: Wrong password!"), http.StatusBadRequest)
		return
	}

	// create a jwt user
	u := jwtUser{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}

	// generate tokens
	tokens, err := app.auth.GenerateTokenPair(&u)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	refreshCookie := app.auth.GetRefreshCookie(tokens.RefreshToken)
	http.SetCookie(w, refreshCookie)

	app.writeJSON(w, http.StatusAccepted, tokens)
}

// refreshToken checks for a valid refresh cookie, and returns a JWT if it finds one.
func (app *application) refreshToken(w http.ResponseWriter, r *http.Request) {
	for _, cookie := range r.Cookies() {
		if cookie.Name == app.auth.CookieName {
			claims := &Claims{}
			refreshToken := cookie.Value

			// parse the token to get the claims
			_, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(app.JWTSecret), nil
			})
			if err != nil {
				app.errorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
				return
			}

			// get the user id from the token claims
			userID, err := strconv.Atoi(claims.Subject)
			if err != nil {
				app.errorJSON(w, errors.New("unknown user"), http.StatusUnauthorized)
				return
			}

			user, err := app.DB.GetUserByID(userID)
			if err != nil {
				app.errorJSON(w, errors.New("unknown user"), http.StatusUnauthorized)
				return
			}

			u := jwtUser{
				ID:        user.ID,
				FirstName: user.FirstName,
				LastName:  user.LastName,
			}

			tokenPairs, err := app.auth.GenerateTokenPair(&u)
			if err != nil {
				app.errorJSON(w, errors.New("error generating tokens"), http.StatusUnauthorized)
				return
			}

			http.SetCookie(w, app.auth.GetRefreshCookie(tokenPairs.RefreshToken))

			app.writeJSON(w, http.StatusOK, tokenPairs)

		}
	}
}

// logout logs the user out by sending an expired cookie to delete the refresh cookie.
func (app *application) logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, app.auth.GetExpiredRefreshCookie())
	w.WriteHeader(http.StatusAccepted)
}

func (app *application) AllThreads(w http.ResponseWriter, r *http.Request) {
	threads, err := app.DB.AllThreads()
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, threads)
}

func (app *application) GetThread(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	threadID, err := strconv.Atoi(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	thread, err := app.DB.SingleThread(threadID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, thread)
}

func (app *application) GetComments(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	threadID, err := strconv.Atoi(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	thread, err := app.DB.GetCommentsByThreadID(threadID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, thread)
}

func (app *application) InsertThread(w http.ResponseWriter, r *http.Request) {
	_, claims, err := app.auth.GetTokenFromHeaderAndVerify(w, r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		app.errorJSON(w, errors.New("unknown user"), http.StatusUnauthorized)
		return
	}

	var thread models.Thread

	err = app.readJSON(w, r, &thread)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	thread.UserID = userID
	thread.CreatedAt = time.Now()
	thread.UpdatedAt = time.Now()

	_, err = app.DB.InsertThread(thread)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

}

func (app *application) InsertComment(w http.ResponseWriter, r *http.Request) {
	_, claims, err := app.auth.GetTokenFromHeaderAndVerify(w, r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		app.errorJSON(w, errors.New("unknown user"), http.StatusUnauthorized)
		return
	}

	var comment models.Comment

	err = app.readJSON(w, r, &comment)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	comment.UserID = userID
	comment.CreatedAt = time.Now()
	comment.UpdatedAt = time.Now()

	_, err = app.DB.InsertComment(comment)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

}

func (app *application) InsertReply(w http.ResponseWriter, r *http.Request) {
	_, claims, err := app.auth.GetTokenFromHeaderAndVerify(w, r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		app.errorJSON(w, errors.New("unknown user"), http.StatusUnauthorized)
		return
	}

	var reply models.Reply

	err = app.readJSON(w, r, &reply)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	reply.UserID = userID
	reply.CreatedAt = time.Now()
	reply.UpdatedAt = time.Now()

	_, err = app.DB.InsertReply(reply)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

}
