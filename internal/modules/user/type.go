package user

import (
	"encoding/json"
	"golang-api-starter/internal/helper"
	"golang-api-starter/internal/modules/groupUser"

	"github.com/golang-jwt/jwt/v4"
)

type UserClaims struct {
	UserId    interface{} `json:"userId"`
	Username  string      `json:"username"`
	TokenType string      `json:"tokenType"`
	jwt.RegisteredClaims
}

// other user related types are grouped in internal/modules/groupUser/user_type.go

type cacheValue struct {
	Users        []*groupUser.User
	Pagination *helper.Pagination
}

// MarshalBinary serializes data into a byte slice for caching.
func (gus *cacheValue) MarshalBinary() (data []byte, err error) {
	bytes, err := json.Marshal(gus)
	return bytes, err
}

// UnmarshalBinary deserializes the byte slice back into data for caching.
func (gus *cacheValue) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, gus)
}

