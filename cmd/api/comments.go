package main

import (
	"net/http"

	"github.com/CP-Payne/social/internal/store"
)

type CreateCommentPayload struct {
	Content string `json:"content" validate:"required,max=1000"`
}

// CreateComment godoc
//
//	@Summary		Creates a comment
//	@Description	Creates a comment
//	@Tags			posts
//	@Param			id	path	int	true	"Post ID"
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		CreateCommentPayload	true	"Comment payload"
//	@Success		201		{object}	store.Comment
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/posts/{id}/comment [post]
func (app *application) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreateCommentPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	post := getPostFromCtx(r)

	user := getUserFromContext(r)

	comment := &store.Comment{
		PostID:  post.ID,
		UserID:  user.ID,
		Content: payload.Content,
	}

	ctx := r.Context()

	if err := app.store.Comments.Create(ctx, comment); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, comment); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
