package auth

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/edgexfoundry/export-go/mongo"
	"go.uber.org/zap"
	"gopkg.in/mgo.v2/bson"
)

func login(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Error("Empty body", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user := &User{}
	if err = json.Unmarshal(body, user); err != nil {
		logger.Error("Malformed JSON", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if user.Username == "" || user.Password == "" {
		logger.Error("Empty field (username or password)", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	s := repo.Session.Copy()
	defer s.Close()
	c := s.DB(mongo.DBName).C(mongo.CollectionName)

	test := &User{}
	if err := c.Find(bson.M{"username": user.Username}).One(&test); err != nil {
		logger.Error("Failed to query by id", zap.Error(err))
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err := CheckPassword(user.Password, test.Password); err != nil {
		logger.Error("Incorrect password", zap.Error(err))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	jwt, err := CreateKey(user.Username)
	if err != nil {
		logger.Error("Failed to create JWT", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, "token:"+jwt)
}
