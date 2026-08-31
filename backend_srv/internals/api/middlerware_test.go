package api

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/require"
	"scm-simple-luke.com/dir/internals/token"
	"scm-simple-luke.com/dir/internals/utils"
)

func addAuthorization(
	t *testing.T,
	req *http.Request,
	username string,
	role string,
	tokenMaker token.TokenMaker,
	durationToken time.Duration,
) (*token.Payload, error) {

	accessToken, payload, err := tokenMaker.CreateToken(username, role, durationToken)
	require.NoError(t, err)
	if err != nil {
		return nil, err
	}

	authBearer := fmt.Sprintf("%s %s", authorizationTypeBearer, accessToken)
	req.Header.Set(authorizationHeaderKey, authBearer)

	// log.Printf("check header: %s", req.Header.Get(authorizationHeaderKey))
	return payload, nil
}

func TestMiddleware(t *testing.T) {

	log.Println("start test")
	authPath := "/auth"
	cfgTest := utils.Config{
		TokenSymmetricKey:   utils.RandomString(32),
		AccessTokenDuration: 200 * time.Millisecond,
	}

	scrt32 := utils.RandomString(32)
	newJWT, err := token.NewJWTMaker(scrt32)
	testServer, err := NewServer(cfgTest, nil, newJWT, nil)
	require.NoError(t, err)
	if err != nil {
		require.NoError(t, err)
	}

	reqTest, err := http.NewRequest("GET", authPath, nil)
	require.NoError(t, err)
	// test error timeOut
	_, err = addAuthorization(t, reqTest, "admin01", "admin", newJWT, 10*time.Millisecond)
	require.NoError(t, err)

	testServer.route.GET(authPath, authMiddlewareV2(newJWT, func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {}))
	recorder := httptest.NewRecorder()
	testServer.route.ServeHTTP(recorder, reqTest)

	log.Printf("status code: %v", recorder.Code)
	require.Equal(t, http.StatusOK, recorder.Code)

}
