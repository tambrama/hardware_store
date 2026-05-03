package auth

import (
	"auth-service/internal/model"
	"auth-service/internal/web/dto"
	"context"
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	pb "github.com/tambrama/protos/gen/go/sso"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serverAPI struct {
	pb.UnimplementedAuthServer
	auth     AuthService
	user     UserService
	token    TokenService
	validate *validator.Validate
	oauth    OAuthService
}

type AuthService interface {
	Login(ctx context.Context, email, password, appID string) (accessToken, refreshToken string, err error)
	Logout(ctx context.Context, userID, appID string) error
}

type UserService interface {
	RegisterNewUser(ctx context.Context, email, password, name, surname, phoneNumber string) (userID string, err error)
	ChangePassword(ctx context.Context, email, oldPassword, newPassword string) error
	RestorePassword(ctx context.Context, email string) error
	GetUserByEmail(ctx context.Context, email string) (model.Users, error)
	GetUserByGoogleID(ctx context.Context, googleID string) (model.Users, error)
	RegisterUserViaGoogle(ctx context.Context, googleID, email, name, surname, phone string) (userID uuid.UUID, err error)
	LinkGoogleToUser(ctx context.Context, userID uuid.UUID, googleID string) error
}

type TokenService interface {
	RefreshToken(ctx context.Context, refreshToken, appId string) (string, string, error)
	ValidateToken(ctx context.Context, token string) (*model.CustomClaims, error)
	IsAdmin(ctx context.Context, userID uuid.UUID) (isAdmin bool, err error)
	IssueTokens(ctx context.Context, email, appID string) (accessToken, refreshToken string, err error)
}

type OAuthService interface {
	GetAuthURL(state string) string
	ExchangeCode(ctx context.Context, code string) (*model.GoogleUserInfo, *oauth2.Token, error)
}

func Register(gRPC *grpc.Server, auth AuthService, user UserService,
	token TokenService, validate *validator.Validate, oauth OAuthService) {
	pb.RegisterAuthServer(gRPC, &serverAPI{
		auth:     auth,
		user:     user,
		token:    token,
		validate: validate,
		oauth:    oauth,
	})
}

func (s *serverAPI) GetGoogleAuthURL(ctx context.Context, req *pb.GetGoogleAuthURLRequest) (*pb.GetGoogleAuthURLResponse, error) {
	if req.GetState() == "" {
		return nil, status.Error(codes.InvalidArgument, "state is required")
	}
	return &pb.GetGoogleAuthURLResponse{
		AuthUrl: s.oauth.GetAuthURL(req.GetState()),
	}, nil
}
func (s *serverAPI) GoogleCallback(ctx context.Context, req *pb.GoogleCallbackRequest) (*pb.LoginResponse, error) {
	// const op = "handler.GoogleCallback"
	requestID := uuid.New().String()[:8]

	fmt.Printf("🔥 [%s] GoogleCallback START\n", requestID)

	if req.GetCode() == "" || req.GetState() == "" || req.GetAppId() == "" {
		return nil, status.Error(codes.InvalidArgument, "code, state and app_id are required")
	}

	googleUser, _, err := s.oauth.ExchangeCode(ctx, req.GetCode())
	if err != nil {
		fmt.Printf("🔥 [%s] ExchangeCode failed: %v\n", requestID, err)
		return nil, status.Error(codes.Internal, "failed to authenticate with Google")
	}

	user, err := s.user.GetUserByGoogleID(ctx, googleUser.GoogleID)
	if err != nil {
		fmt.Printf("🔥 [%s] GetUserByGoogleID: %v\n", requestID, err)
		if errors.Is(err, model.ErrUserNotFound) {
			fmt.Printf("🔥 [%s] Searching by email: %s\n", requestID, googleUser.Mail)
			user, err = s.user.GetUserByEmail(ctx, googleUser.Mail)
			if err != nil {
				fmt.Printf("🔥 [%s] GetUserByEmail: %v\n", requestID, err)
				if errors.Is(err, model.ErrUserNotFound) {
					fmt.Printf("🔥 [%s] Registering new user via Google\n", requestID)
					userID, err := s.user.RegisterUserViaGoogle(ctx, googleUser.GoogleID, googleUser.Mail, googleUser.Name, googleUser.Surname, "")
					if err != nil {
						fmt.Printf("🔥 [%s] RegisterUserViaGoogle FAILED: %v\n", requestID, err)
						return nil, status.Error(codes.Internal, "failed to create user")
					}
					fmt.Printf("🔥 [%s] User registered: %s\n", requestID, userID)
					user = model.Users{ID: userID, Mail: googleUser.Mail}
				} else {
					return nil, status.Error(codes.Internal, "failed to get user by email")
				}
			} else {
				fmt.Printf("🔥 [%s] Linking GoogleID to existing user\n", requestID)
				err = s.user.LinkGoogleToUser(ctx, user.ID, googleUser.GoogleID)
				if err != nil {
					fmt.Printf("🔥 [%s] LinkGoogleToUser FAILED: %v\n", requestID, err)
					return nil, status.Error(codes.Internal, "failed to link Google to user")
				}
			}
		} else {
			return nil, status.Error(codes.Internal, "failed to get user by GoogleID")
		}
	}

	fmt.Printf("🔥 [%s] Checking password for user: %s\n", requestID, user.Mail)

	if user.Password == nil || len(*user.Password) == 0 {
		fmt.Printf("🔥 [%s] User has no password, issuing tokens via Google\n", requestID)
		accessToken, refreshToken, err := s.token.IssueTokens(ctx, user.Mail, req.GetAppId())
		if err != nil {
			fmt.Printf("🔥 [%s] IssueTokens FAILED: %v\n", requestID, err)
			return nil, status.Error(codes.Internal, "failed to issue tokens")
		}
		fmt.Printf("🔥 [%s] Tokens issued successfully\n", requestID)
		return &pb.LoginResponse{AccessToken: accessToken, RefreshToken: refreshToken}, nil
	}

	fmt.Printf("🔥 [%s] Calling Login with password\n", requestID)
	accessToken, refreshToken, err := s.auth.Login(ctx, user.Mail, "", req.GetAppId())
	if err != nil {
		fmt.Printf("🔥 [%s] Login FAILED: %v\n", requestID, err)
		return nil, status.Error(codes.Internal, "failed to login")
	}

	fmt.Printf("🔥 [%s] GoogleCallback COMPLETE\n", requestID)
	return &pb.LoginResponse{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (s *serverAPI) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	input := &dto.LoginInput{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
		AppID:    req.GetAppId(),
	}

	if err := s.validate.Struct(input); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			field := ve[0].Field()
			switch field {
			case "Email":
				return nil, status.Error(codes.InvalidArgument, "empty or invalid email")
			case "Password":
				return nil, status.Error(codes.InvalidArgument, "password is required")
			case "AppID":
				return nil, status.Error(codes.InvalidArgument, "app_id is required")
			default:
				return nil, status.Errorf(codes.InvalidArgument, "invalid field: %s", field)
			}
		}
		return nil, status.Error(codes.InvalidArgument, "invalid input")
	}

	accessToken, refreshToken, err := s.auth.Login(ctx, input.Email, input.Password, input.AppID)
	if err != nil {
		if errors.Is(err, model.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid email or password")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &pb.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *serverAPI) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	input := &dto.RegisterInput{
		Email:       req.GetEmail(),
		Password:    req.GetPassword(),
		Name:        req.GetName(),
		Surname:     req.GetSurname(),
		PhoneNumber: req.GetPhoneNumber(),
	}

	if err := s.validate.Struct(input); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			field := ve[0].Field()

			switch field {
			case "Email":
				return nil, status.Error(codes.InvalidArgument, "invalid or empty email")
			case "Password":
				return nil, status.Error(codes.InvalidArgument, "password is required")
			case "Name":
				return nil, status.Error(codes.InvalidArgument, "name is required")
			case "Surname":
				return nil, status.Error(codes.InvalidArgument, "surname is required")
			case "PhoneNumber":
				return nil, status.Error(codes.InvalidArgument, "invalid phone number")
			default:
				return nil, status.Errorf(codes.InvalidArgument, "invalid field: %s", field)
			}
		}
		return nil, status.Error(codes.InvalidArgument, "invalid input")
	}
	userID, err := s.user.RegisterNewUser(ctx, input.Email, input.Password, input.Name, input.Surname, input.PhoneNumber)
	if err != nil {
		if errors.Is(err, model.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &pb.RegisterResponse{
		UserId: userID,
	}, nil
}

func (s *serverAPI) IsAdmin(ctx context.Context, req *pb.IsAdminRequest) (*pb.IsAdminResponse, error) {
	if req.GetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	isAdmin, err := s.token.IsAdmin(ctx, uuid.MustParse(req.GetUserId()))
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &pb.IsAdminResponse{
		IsAdmin: isAdmin,
	}, nil
}

func (s *serverAPI) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*pb.Empty, error) {
	input := &dto.ChangePasswordInput{
		Email:       req.GetEmail(),
		OldPassword: req.GetOldPassword(),
		NewPassword: req.GetNewPassword(),
	}
	if err := s.validate.Struct(input); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid input")
	}
	err := s.user.ChangePassword(ctx, req.GetEmail(), input.OldPassword, input.NewPassword)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &pb.Empty{}, nil
}

func (s *serverAPI) RestorePassword(ctx context.Context, req *pb.RestorePasswordRequest) (*pb.Empty, error) {
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}
	err := s.user.RestorePassword(ctx, req.Email)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &pb.Empty{}, nil
}

func (s *serverAPI) Refresh(ctx context.Context, req *pb.RefreshRequest) (*pb.LoginResponse, error) {
	input := &dto.RefreshInput{
		RefreshToken: req.GetRefreshToken(),
		AppID:        req.GetAppId(),
	}
	if err := s.validate.Struct(input); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid input")
	}

	access, refresh, err := s.token.RefreshToken(ctx, input.RefreshToken, input.AppID)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &pb.LoginResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}
func (s *serverAPI) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.Empty, error) {
	input := &dto.LogoutInput{
		UserID: req.GetUserId(),
		AppID:  req.GetAppId(),
	}
	if input.UserID == "" || input.AppID == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id and app_id are required")
	}
	if err := s.auth.Logout(ctx, input.UserID, input.AppID); err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &pb.Empty{}, nil
}

func (s *serverAPI) Validate(ctx context.Context, req *pb.ValidateRequest) (*pb.ValidateResponse, error) {
	if req.GetToken() == "" {
		return nil, status.Error(codes.InvalidArgument, "token is required")
	}
	claims, err := s.token.ValidateToken(ctx, req.GetToken())
	if err != nil {
		if errors.Is(err, model.ErrInvalidCredentials) || errors.Is(err, model.ErrInvalidToken) {
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &pb.ValidateResponse{
		UserId: claims.UserID.String(),
		AppId:  claims.AppID.String(),
	}, nil
}
