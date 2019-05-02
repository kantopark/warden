package application

import (
	"io/ioutil"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/jwtauth"
	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"warden/utils"
)

var (
	tokenDuration time.Duration
	jwtToken      *jwtauth.JWTAuth
)

func authInit() {
	alg := utils.StrUpperTrim(viper.GetString("auth.signing_alg"))
	if alg == "" {
		log.Fatalln("JWT signing algorithm must be specified")
	}

	pubKey := viper.GetString("auth.public_key")
	privKey := viper.GetString("auth.private_key")

	var privateKey interface{} // used to sign
	var publicKey interface{}  // used to verify

	switch alg {
	case "HS256", "HS384", "HS512":
		privateKey = []byte(privKey + pubKey)
		publicKey = nil
	case "RS256", "RS384", "RS512":
		bytePrivKey, err := ioutil.ReadFile(privKey)
		fatalIfError(err)

		privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(bytePrivKey)
		fatalIfError(err)

		bytePubKey, err := ioutil.ReadFile(pubKey)
		fatalIfError(err)

		publicKey, err = jwt.ParseRSAPublicKeyFromPEM(bytePubKey)
		fatalIfError(err)

	case "ES256", "ES384", "ES512":
		bytePrivKey, err := ioutil.ReadFile(privKey)
		fatalIfError(err)

		privateKey, err = jwt.ParseECPrivateKeyFromPEM(bytePrivKey)
		fatalIfError(err)

		bytePubKey, err := ioutil.ReadFile(pubKey)
		fatalIfError(err)

		publicKey, err = jwt.ParseECPublicKeyFromPEM(bytePubKey)
		fatalIfError(err)
	default:
		log.Fatalf("unknown JWT algorithm: %s\n", alg)
	}

	jwtToken = jwtauth.New(alg, privateKey, publicKey)
	if dur, err := time.ParseDuration(viper.GetString("auth.expiry")); err != nil {
		tokenDuration = 1 * time.Hour
	} else {
		tokenDuration = dur
	}
}

func createToken(username, email string) (string, error) {
	token := jwtauth.Claims(jwt.MapClaims{
		"username": username,
		"email":    email,
	}).SetExpiryIn(tokenDuration)

	_, tokenString, err := jwtToken.Encode(token)
	if err != nil {
		return "", errors.Wrap(err, "error tokenizing claim")
	}
	return tokenString, nil
}
