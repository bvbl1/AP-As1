package grpc

import (
	"Assignment1_AbylayMoldakhmet/proto/gen"
	"Assignment1_AbylayMoldakhmet/user-service/internal/domain"
	"Assignment1_AbylayMoldakhmet/user-service/internal/usecase"
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserServer struct {
	gen.UnimplementedUserServiceServer
	authUC usecase.AuthUsecase
	userUC usecase.UserUsecase
}

func NewUserServer(authUC usecase.AuthUsecase, userUC usecase.UserUsecase) *UserServer {
	return &UserServer{authUC: authUC, userUC: userUC}
}

func (s *UserServer) Register(ctx context.Context, req *gen.RegisterRequest) (*gen.UserResponse, error) {
	user, err := s.authUC.Register(req.Email, req.Password)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return convertUserToResponse(user), nil
}

func (s *UserServer) Login(ctx context.Context, req *gen.LoginRequest) (*gen.LoginResponse, error) {
	token, err := s.authUC.Login(req.Email, req.Password)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}

	return &gen.LoginResponse{
		AccessToken: token,
	}, nil
}

func (s *UserServer) GetUserProfile(ctx context.Context, req *gen.UserIDRequest) (*gen.UserResponse, error) {
	user, err := s.userUC.GetByID(req.UserId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	return convertUserToResponse(user), nil
}

func (s *UserServer) UpdateUser(ctx context.Context, req *gen.UpdateUserRequest) (*gen.UserResponse, error) {
	objID, err := primitive.ObjectIDFromHex(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user ID")
	}

	user := &domain.User{
		ID:    objID,
		Email: req.Email,
		Role:  req.Role,
	}

	if err := s.userUC.Update(user); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return convertUserToResponse(user), nil
}

func (s *UserServer) DeleteUser(ctx context.Context, req *gen.UserIDRequest) (*emptypb.Empty, error) {
	if err := s.userUC.Delete(req.UserId); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

// Вспомогательная функция для конвертации
func convertUserToResponse(user *domain.User) *gen.UserResponse {
	return &gen.UserResponse{
		Id:    user.ID.Hex(),
		Email: user.Email,
		Role:  user.Role,
	}
}
