package auth

import (
	"auth-service/internal/model"
	"auth-service/internal/web/dto"
	"context"
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	pb "github.com/tambrama/protos/gen/go/sso"
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
}

type AuthService interface {
	Login(ctx context.Context, email, password, appID string) (accessToken, refreshToken string, err error)
	Logout(ctx context.Context, token string) error
}

type UserService interface {
	RegisterNewUser(ctx context.Context, email, password, name, surname, phoneNumber string) (userID string, err error)
	ChangePassword(ctx context.Context, email, oldPassword, newPassword string) error
	RestorePassword(ctx context.Context, email string) error
}

type TokenService interface {
	RefreshToken(ctx context.Context, refreshToken, appId string) (string, string, error)
	ValidateToken(ctx context.Context, token string) (*model.CustomClaims, error)
	IsAdmin(ctx context.Context, userID uuid.UUID) (isAdmin bool, err error)
}

func Register(gRPC *grpc.Server, auth AuthService, user UserService,
	token TokenService, validate *validator.Validate) {
	pb.RegisterAuthServer(gRPC, &serverAPI{
		auth:     auth,
		user: user,
		token: token,
		validate: validate,
	})
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
		RefreshToken: req.GetRefreshToken(),
	}
	if input.RefreshToken == "" {
		return nil, status.Error(codes.InvalidArgument, "refresh token is required")
	}
	if err := s.auth.Logout(ctx, input.RefreshToken); err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &pb.Empty{}, nil
}

func (s *serverAPI) Validate(ctx context.Context, req *pb.ValidateRequest) (*pb.ValidateResponse, error) {
	if req.GetToken() == "" {
		return nil, status.Error(codes.InvalidArgument, "token is required")
	}
	claims, err :=s.token.ValidateToken(ctx, req.GetToken())
	if err != nil {
		if errors.Is(err, model.ErrInvalidCredentials) {
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &pb.ValidateResponse{
		UserId: claims.UserID.String(),
		AppId:  claims.AppID.String(),
	}, nil
}