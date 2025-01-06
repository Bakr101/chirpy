package auth

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestHashing(t *testing.T) {
	hashed_password, err:= HashPassword("bakr1")
	if err != nil {
		fmt.Printf("hashed_pass: %v", hashed_password)
		t.Errorf("error: %v", err)
	}
	//fmt.Printf("hashed password: %v", hashed_password)
}

func TestHashAndCompare(t *testing.T) {
	//fail case flip either HashPassword or checkPasswordHash inputs
	//hashed_password, err := HashPassword("")
	//success case
	hashed_password, err := HashPassword("bakr1")
	if err != nil {
		t.Errorf("error hashing password: %v", err)
	}
	
	//err = CheckPasswordHash("", hashed_password)
	//success case
	err = CheckPasswordHash("bakr1", hashed_password)
	if err != nil {
		t.Errorf("err in Check password hash: %v", err)
	}
	//fmt.Printf("Hashed pasword from check: %v, err: %v",hashed_password, err)
}

func TestGenerateToken(t *testing.T){
	_, err := MakeJWT(uuid.New(), "bakr", 10 * time.Minute)
	if err != nil {
		t.Errorf("error generating token: %v", err)
	}
	//fmt.Println(stringifedToken)
}

func TestValidateToken(t *testing.T){
	validUUID := uuid.New()
	
	fmt.Printf("valid ID: %v\n",validUUID)
	stringifiedToken, err := MakeJWT(validUUID, "bakr", 10)
	if err != nil {
		t.Errorf("%v", err)
	}
	result, err := ValidateJWT(stringifiedToken, "bakr")
	fmt.Printf("Result: %v\n", result)
	if err != nil {
		t.Errorf("%v", err)
	}
	

}

func TestGetBearerToken(t *testing.T){ 
	headers:= http.Header{
		"Authorization": []string{"Bearer TOKEN_STRING"},
	}
	token, err := GetBearerToken(headers)
	t.Log(token)
	if err != nil{
		t.Errorf("%v", err)
	}
}

func TestMakreRefreshToken(t *testing.T){
	token, err := MakeRefreshToken()
	if err != nil {
		t.Errorf("error making token: %v", err)
	}
	t.Log(token)
}