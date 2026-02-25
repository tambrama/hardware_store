package auth

import (
	"auth-service/internal/web/dto"
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	pb "github.com/tambrama/protos/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serverAPI struct {
	pb.UnimplementedAuthServer
	auth     Auth
	validate *validator.Validate
}

type Auth interface {
	Login(ctx context.Context, email, password, appID string) (accessToken, refreshToken string, err error)
	RegisterNewUser(ctx context.Context, email, password, name, surname, phoneNumber string) (userID string, err error)
	IsAdmin(ctx context.Context, userID uuid.UUID) (isAdmin bool, err error)
	ChangePassword(ctx context.Context, email, oldPassword, newPassword string) error
	RestorePassword(ctx context.Context, email string) error
}

func Register(gRPC *grpc.Server, auth Auth, validate *validator.Validate) {
	pb.RegisterAuthServer(gRPC, &serverAPI{
		auth:     auth,
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
		return nil, status.Error(codes.InvalidArgument, "invalid input")
	}

	accessToken, refreshToken, err := s.auth.Login(ctx, input.Email, input.Password, input.AppID)
	if err != nil {
		/////
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
		return nil, status.Error(codes.InvalidArgument, "invalid input")
	}
	userID, err := s.auth.RegisterNewUser(ctx, input.Email, input.Password, input.Name, input.Surname, input.PhoneNumber)
	if err != nil {
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

	isAdmin, err := s.auth.IsAdmin(ctx, uuid.MustParse(req.GetUserId()))
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
	err := s.auth.ChangePassword(ctx, req.GetEmail(), input.OldPassword, input.NewPassword)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &pb.Empty{}, nil
}

func (s *serverAPI) RestorePassword(ctx context.Context, req *pb.RestorePasswordRequest) (*pb.Empty, error) {
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}
	err := s.auth.RestorePassword(ctx, req.Email)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &pb.Empty{}, nil
}
