package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/go-zoo/bone"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"

	"github.com/edgexfoundry/export-go/mongo"
)

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// CreateUser creates new user account based on provided username and password.
// The account is assigned with one master key - a key with all permissions on
// all owned resources regardless of their type. Provided password in encrypted
// using bcrypt algorithm.
func CreateUser(username, password string) (User, error) {
	u, err := uuid.NewV4()
	if err != nil {
		logger.Error("Failed to generate uuid", zap.Error(err))
		return User{}, err
	}

	user := User{
		ID:       u.String(),
		Username: username,
	}

	p, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("Error encrypting user password", zap.Error(err))
		return user, &AuthError{Code: http.StatusInternalServerError}
	}

	user.Password = string(p)

	return user, nil
}

// CheckPassword tries to determine whether or not the submitted password
// matches the one stored (and hashed) during registration. An error will be
// used to indicate an invalid password.
func CheckPassword(plain, hashed string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
}

func createUser(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Error("Failed to read body for createUser()", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	data := &User{}
	if err = json.Unmarshal(body, data); err != nil {
		logger.Error("Failed to unmarshal JSON", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if data.Username == "" || data.Password == "" {
		logger.Error("Empty field (username or password)", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user, err := CreateUser(data.Username, data.Password)
	if err != nil {
		logger.Error("Error creating user struct", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	s := repo.Session.Copy()
	defer s.Close()
	c := s.DB(mongo.DBName).C(mongo.CollectionName)

	count, err := c.Find(bson.M{"name": user.Username}).Count()
	if err != nil {
		logger.Error("Failed to add user", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if count != 0 {
		logger.Error("Username already taken: " + user.Username)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := c.Insert(user); err != nil {
		logger.Error("Failed to insert user", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("/users/%s", user.ID))
	w.WriteHeader(http.StatusCreated)
}

func getAllUsers(w http.ResponseWriter, r *http.Request) {
	s := repo.Session.Copy()
	defer s.Close()
	c := s.DB(mongo.DBName).C(mongo.CollectionName)

	users := []User{}
	if err := c.Find(nil).All(&users); err != nil {
		logger.Error("Failed to query all users", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(users)
	if err != nil {
		logger.Error("Failed to query all users", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, string(res))
}

func getUserByID(w http.ResponseWriter, r *http.Request) {
	id := bone.GetValue(r, "id")

	s := repo.Session.Copy()
	defer s.Close()
	c := s.DB(mongo.DBName).C(mongo.CollectionName)

	user := User{}
	if err := c.Find(bson.M{"id": id}).One(&user); err != nil {
		logger.Error("Failed to query by id", zap.Error(err))
		w.WriteHeader(http.StatusNotFound)
		return
	}

	res, err := json.Marshal(user)
	if err != nil {
		logger.Error("Failed to query by id", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, string(res))
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	id := bone.GetValue(r, "id")

	s := repo.Session.Copy()
	defer s.Close()
	c := s.DB(mongo.DBName).C(mongo.CollectionName)

	if err := c.Remove(bson.M{"id": id}); err != nil {
		logger.Error("Failed to delete user", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
