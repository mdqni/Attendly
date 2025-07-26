package grpc

import (
	"context"
	authv1 "github.com/mdqni/Attendly/proto/gen/go/auth/v1"
	"github.com/mdqni/Attendly/services/auth/internal/service"
	"github.com/mdqni/Attendly/shared/errs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"strings"
)

type authServer struct {
	authv1.UnimplementedAuthServiceServer
	service service.AuthService
}

func Register(gRPCServer *grpc.Server, svc service.AuthService) {
	authv1.RegisterAuthServiceServer(gRPCServer, &authServer{service: svc})
}

func (h *authServer) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.AuthResponse, error) {
	if strings.TrimSpace(req.Barcode) == "" || strings.TrimSpace(req.Password) == "" {
		return nil, status.Error(codes.InvalidArgument, errs.ErrMissingField.Error())
	}

	resp, err := h.service.Login(ctx, service.LoginInput{
		Barcode:  req.GetBarcode(),
		Password: req.GetPassword(),
	})
	if err != nil {
		log.Println("login error:", err)
		return nil, err
	}

	return &authv1.AuthResponse{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
	}, nil
}

func (h *authServer) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.AuthResponse, error) {
	const op = "grpc.auth.register"
	log.Println("op", op)
	if strings.TrimSpace(req.Name) == "" ||
		strings.TrimSpace(req.Barcode) == "" ||
		strings.TrimSpace(req.Password) == "" ||
		strings.TrimSpace(req.Role) == "" {
		return nil, status.Error(codes.InvalidArgument, errs.ErrMissingField.Error())
	}

	resp, err := h.service.Register(ctx, service.RegisterInput{
		Name:     req.GetName(),
		Barcode:  req.GetBarcode(),
		Password: req.GetPassword(),
		Email:    req.GetEmail(),
		Role:     req.GetRole(),
	})
	log.Println("Resp: ", resp)
	if err != nil {
		log.Println("register error:", err)
		return nil, err
	}

	return &authv1.AuthResponse{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		User:         &resp.User,
	}, nil
}
func (h *authServer) GetUserInfoById(ctx context.Context, request *authv1.GetUserInfoRequest) (*authv1.GetUserInfoResponse, error) {
	user, err := h.service.GetUserInfoById(ctx, request.UserId)
	if err != nil {
		return nil, err
	}
	return &authv1.GetUserInfoResponse{
		Name:    user.Name,
		Barcode: user.Barcode,
		Role:    user.Role,
	}, nil
}

func (h *authServer) Refresh(ctx context.Context, req *authv1.RefreshTokenRequest) (*authv1.AuthResponse, error) {
	refresh, err := h.service.Refresh(ctx, req)
	if err != nil {
		return nil, err
	}
	return &authv1.AuthResponse{
		AccessToken:  refresh.AccessToken,
		RefreshToken: refresh.RefreshToken,
	}, nil
}
