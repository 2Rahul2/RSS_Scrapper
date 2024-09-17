package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/2Rahul2/rssagg/internal/database"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func (apiCfg apiConfig) handlerCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}
	decoder := json.NewDecoder(r.Body)
	var params parameters
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Err: Could not parse follow feed from json : %v", err))
		return
	}
	feed, err := apiCfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    params.FeedID,
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Err: could not create feed follow :%v", err))
		return
	}
	respondWithJSON(w, 201, feed)
}

func (apiCfg apiConfig) handlerGetFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollow, err := apiCfg.DB.GetFeedFollows(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Err: Could not Get follow feed :%v", err))
		return
	}
	respondWithJSON(w, 200, feedFollow)
}

func (apiCfg apiConfig) handlerDeleteFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	var idStr string = chi.URLParam(r, "feedFollowId")
	feedFollowId, err := uuid.Parse(idStr)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("ERR: Could not parse the ID : %v", err))
		return
	}
	err = apiCfg.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		ID:     feedFollowId,
		UserID: user.ID,
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("ERR: Could not delete the feed : %v", err))
		return
	}
	respondWithJSON(w, 200, struct{}{})

}
