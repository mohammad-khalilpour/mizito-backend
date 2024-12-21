package repositories

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"mizito/pkg/models"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)


const secret = ""


type AuthDB interface {
    AuthenticateUser(username string, password string) (bool, error)
    StoreUser(user *models.User) (uint, error)
}


type AuthHandler struct {
	DB AuthDB
}


func (ah *AuthHandler) authenticate (username string, password string) error{
    hashPass := sha256.Sum256([]byte(password))
    hashPassString := hex.EncodeToString(hashPass[:])

    found, err := ah.DB.AuthenticateUser(username, hashPassString)
    if !found {
        return fmt.Errorf("no such username was found: %s", username)
    }
    if err != nil {
        return fmt.Errorf("authentication has failed, err : %s", err)
    }

    return nil  
}

func (ah *AuthHandler) Basic(token string) error{
    rawTkn, err := base64.StdEncoding.DecodeString(token)
    if err != nil {
        return fmt.Errorf("failed to parse and decode token, err : %s", err.Error())
    }
    if tokens := strings.Split(string(rawTkn), ":"); len(tokens) != 2 {
        return fmt.Errorf("failed to parse Auth token, token must contain username and password separated by comma")
    } else {
        return ah.authenticate(tokens[0], tokens[1])
    }
}

func (ah *AuthHandler) Bearer(token string) error {

    jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
        if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected method: %s", t.Header["alg"])
        }
        return secret, nil
    })

    // complete this section for further permission check
    return nil
}

func (ah *AuthHandler) StoreUser(user *models.User) (uint, error) {
    return ah.DB.StoreUser(user)
}