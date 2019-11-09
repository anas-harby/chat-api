package handlers

import (
	"log"
	"time"
	"strconv"
	"strings"
	"net/http"
	"encoding/json"
	"github.com/imroc/req"
	"github.com/gorilla/mux"
	"github.com/anas-harby/go-chat-creation-api/pkg/sidekiq"
	"github.com/anas-harby/go-chat-creation-api/pkg/redis"
	"github.com/anas-harby/go-chat-creation-api/pkg/network"
	"github.com/anas-harby/go-chat-creation-api/configs"
)

type chatResponse struct {
	Number        int64   `json:"number"`
	AccessToken   string  `json:"access_token"`
}

type appsApiResponse struct {
	Name          string  `json:"name"`
	AccessToken   string  `json:"access_token"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
	ChatCount     int64   `json:"chat_count"`
}

func CreateChat(w http.ResponseWriter, r *http.Request) {
	// Read in request
	accessToken := mux.Vars(r)["access_token"]

	// Get next number
	redisClient, err := redis.GetRedis()
	if err != nil {
		network.RespondErr(w, err)
		return
	}

	redisLocker := redis.GetLocker()
	key := configs.RedisChatKeyPrefix + accessToken

	// Begin critical section
	lock, err := redisLocker.Obtain(key + "_LOCK", 1000*time.Millisecond, nil)
	if err != nil {
		defer lock.Release()
		network.RespondErr(w, err)
		return
	}

	exists, err := redisClient.Exists(key).Result()
	if err != nil {
		defer lock.Release()
		network.RespondErr(w, err)
		return
	} else if exists == 0 {
		log.Println("Key " + key + " not found in Redis, requsting chat count from Rails instead")
		appsResp, err := RequestChats(accessToken)
		if err != nil {
			defer lock.Release()
			network.RespondErr(w, err)
			return
		}
		redisClient.Set(key, appsResp.ChatCount, 1)
	}

	nextNum, err := redisClient.Incr(key).Result()
	defer lock.Release()
	// End critical section
	if err != nil {
		network.RespondErr(w, err)
		return
	}

	// Push to Sidekiq
	err = sidekiq.Push(configs.RedisChatQueue, configs.ChatWorkerClass, accessToken, strconv.FormatInt(nextNum, 10))
	if err != nil {
		network.RespondErr(w, err)
		return
	}

	// Send response
	resp := chatResponse{Number: nextNum, AccessToken: accessToken}
	respBytes, err := json.Marshal(resp)
	if err != nil {
		network.RespondErr(w, err)
		return
	}

	network.Respond(w, respBytes, http.StatusCreated)
}

func RequestChats(accessToken string) (appsApiResponse, error) {
	var resp appsApiResponse

	r, err := req.Get(strings.Replace(configs.AppEndpoint + configs.ChatsRoute, "{access_token}", accessToken, 1))
	if err != nil {
		return resp, err
	}

	r.ToJSON(&resp)
	return resp, nil
}

