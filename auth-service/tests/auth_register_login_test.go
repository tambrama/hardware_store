package tests

import (
	"auth-service/tests/suite"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	pb "github.com/tambrama/protos/gen/go/sso"
)

const (
	emptyAppId     = ""
	appId          = "00000000-0000-0000-0000-000000000001"
	passDefaultLen = 10
)

func TestRegisterLogin_Login_HappyPath(t *testing.T) {
	ctx, st := suite.NewSuite(t)

	email := gofakeit.Email()
	password := randomPassword()
	name := gofakeit.FirstName()
	surname := gofakeit.LastName()
	phone := gofakeit.Phone()
	respReq, err := st.AuthClient.Register(ctx, &pb.RegisterRequest{
		Email:       email,
		Password:    password,
		Name:        name,
		Surname:     surname,
		PhoneNumber: phone,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, respReq.GetUserId())
	loginTime := time.Now()
	respLogin, err := st.AuthClient.Login(ctx, &pb.LoginRequest{
		Email:    email,
		Password: password,
		AppId:    appId,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, respLogin.GetAccessToken())
	assert.NotEmpty(t, respLogin.GetRefreshToken())

	token := respLogin.GetAccessToken()

	tokenParsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(st.Cfg.JWTSecretKey), nil
	})

	require.NoError(t, err)

	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	require.True(t, ok)
	assert.Equal(t, respReq.GetUserId(), claims["user_id"])
	assert.Equal(t, email, claims["email"].(string))
	assert.Equal(t, appId, claims["app_id"].(string))

	assert.True(t, tokenParsed.Valid)

	const expectedExp = 1

	assert.InDelta(t, loginTime.Add(st.Cfg.TokenTTL).Unix(), claims["exp"].(float64), expectedExp)

}

func TestRegisterLogin_DuplicatedRegistration(t *testing.T) {
	ctx, st := suite.NewSuite(t)

	email := gofakeit.Email()
	pass := randomPassword()
	name := gofakeit.FirstName()
	surname := gofakeit.LastName()
	phone := gofakeit.Phone()

	// Первая попытка должна быть успешной
	respReg, err := st.AuthClient.Register(ctx, &pb.RegisterRequest{
		Email:       email,
		Password:    pass,
		Name:        name,
		Surname:     surname,
		PhoneNumber: phone,
	})
	require.NoError(t, err)
	require.NotEmpty(t, respReg.GetUserId())

	// Вторая попытка - фэил
	respReg, err = st.AuthClient.Register(ctx, &pb.RegisterRequest{
		Email:       email,
		Password:    pass,
		Name:        name,
		Surname:     surname,
		PhoneNumber: phone,
	})
	require.Error(t, err)
	assert.Empty(t, respReg.GetUserId())
	assert.ErrorContains(t, err, "user already exists")
}

func TestRegister_FailCases(t *testing.T) {
	ctx, st := suite.NewSuite(t)

	tests := []struct {
		name        string
		email       string
		password    string
		expectedErr string
		userName    string
		surname     string
		phone       string
	}{
		{
			name:        "Register with Empty Password",
			email:       gofakeit.Email(),
			password:    "",
			expectedErr: "password is required",
			userName:    gofakeit.FirstName(),
			surname:     gofakeit.LastName(),
			phone:       gofakeit.Phone(),
		},
		{
			name:        "Register with Empty Email",
			email:       "",
			password:    randomPassword(),
			expectedErr: "invalid or empty email",
			userName:    gofakeit.FirstName(),
			surname:     gofakeit.LastName(),
			phone:       gofakeit.Phone(),
		},
		{
			name:        "Register with Both Empty",
			email:       "",
			password:    "",
			expectedErr: "invalid or empty email",
			userName:    gofakeit.FirstName(),
			surname:     gofakeit.LastName(),
			phone:       gofakeit.Phone(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := st.AuthClient.Register(ctx, &pb.RegisterRequest{
				Email:       tt.email,
				Password:    tt.password,
				Name:        tt.userName,
				Surname:     tt.surname,
				PhoneNumber: tt.phone,
			})
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.expectedErr)

		})
	}
}
func TestLogin_FailCases(t *testing.T) {
	ctx, st := suite.NewSuite(t)

	tests := []struct {
		name        string
		email       string
		password    string
		appID       string
		expectedErr string
	}{
		{
			name:        "Login with Empty Password",
			email:       gofakeit.Email(),
			password:    "",
			appID:       appId,
			expectedErr: "password is required",
		},
		{
			name:        "Login with Empty Email",
			email:       "",
			password:    randomPassword(),
			appID:       appId,
			expectedErr: "empty or invalid email",
		},
		{
			name:        "Login with Both Empty Email and Password",
			email:       "",
			password:    "",
			appID:       appId,
			expectedErr: "empty or invalid email",
		},
		{
			name:        "Login with Non-Matching Password",
			email:       gofakeit.Email(),
			password:    randomPassword(),
			appID:       appId,
			expectedErr: "invalid email or password",
		},
		{
			name:        "Login without AppID",
			email:       gofakeit.Email(),
			password:    randomPassword(),
			appID:       emptyAppId,
			expectedErr: "app_id is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := st.AuthClient.Register(ctx, &pb.RegisterRequest{
				Email:       gofakeit.Email(),
				Password:    randomPassword(),
				Name:        gofakeit.FirstName(),
				Surname:     gofakeit.LastName(),
				PhoneNumber: gofakeit.Phone(),
			})
			require.NoError(t, err)

			_, err = st.AuthClient.Login(ctx, &pb.LoginRequest{
				Email:    tt.email,
				Password: tt.password,
				AppId:    tt.appID,
			})
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.expectedErr)
		})
	}
}

func randomPassword() string {
	return gofakeit.Password(true, true, true, true, false, passDefaultLen)
}
