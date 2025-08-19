package services

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/sidiik/moonpay/auth_service/internal/constants"
	"github.com/sidiik/moonpay/auth_service/internal/domain"
	"github.com/sidiik/moonpay/auth_service/internal/infra/config"
	"github.com/sidiik/moonpay/auth_service/internal/infra/rabbitmq"
	"github.com/sidiik/moonpay/auth_service/internal/utils"
	authpb "github.com/sidiik/moonpay/auth_service/proto"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OtpService struct {
	otpRepo  domain.OtpRepository
	userRepo domain.UserRepository
	rabbit   *rabbitmq.RabbitMQ
	log      domain.Logger
}

func NewOtpService(otpRepo domain.OtpRepository, userRepo domain.UserRepository, rabbit *rabbitmq.RabbitMQ, log domain.Logger) *OtpService {
	return &OtpService{
		otpRepo:  otpRepo,
		userRepo: userRepo,
		log:      log,
		rabbit:   rabbit,
	}
}

func (s *OtpService) RequestPasswordReset(ctx context.Context, req *authpb.RequestPasswordResetRequest) error {
	s.log.Info("Check if the user exists")
	formattedEmail := strings.ToLower(req.Email)
	user, err := s.userRepo.GetUserByEmail(ctx, formattedEmail)

	if err != nil {
		return status.Error(codes.NotFound, constants.ErrUserNotFound)
	}

	s.log.Info("Check if there is any active otp")
	otp, _ := s.otpRepo.CheckActiveOtp(ctx, user.ID)
	if otp != nil {
		s.log.Error("failed to get active otp", "error", err)
		return status.Error(codes.Internal, constants.ErrOtpResendLimit)
	}

	otpCode, err := utils.GenerateOtpCode()
	if err != nil {
		s.log.Error("failed to generate otp code", "error", err)
		return status.Error(codes.Internal, constants.ErrOtpGenerationFail)
	}

	otpExpire, _ := strconv.ParseInt(config.AppConfig.OtpCodeExpire, 10, 64)

	newOtp := domain.Otp{
		Code:         otpCode,
		ExpiresAt:    time.Now().Add(time.Minute * time.Duration(otpExpire)).UTC(),
		NextResendAt: time.Now().Add(time.Minute * 2).UTC(),
		OtpReason:    domain.PasswordReset,
		UserID:       user.ID,
	}

	if err := s.otpRepo.CreateOtp(ctx, &newOtp); err != nil {
		s.log.Error("failed to create otp", "error", err)
		return status.Error(codes.Internal, constants.ErrOtpGenerationFail)
	}

	event := map[string]any{
		"event":     "user.otp.requestPasswordReset",
		"email":     user.Email,
		"fullName":  user.FullName,
		"code":      otpCode,
		"expiresIn": config.AppConfig.OtpCodeExpire,
	}

	body, err := json.Marshal(event)
	if err != nil {
		s.log.Error("failed to parse otp event", "error", err)
		return status.Error(codes.Internal, constants.ErrOtpGenerationFail)
	}

	s.log.Info("Publish password reset request event")
	if err := s.rabbit.Publish("user.otp.requestPasswordReset", body); err != nil {
		s.log.Error("failed to publish password reset request event", "error", err)
		return status.Error(codes.Internal, constants.ErrOtpGenerationFail)
	}

	return nil

}

func (s *OtpService) ResetPassword(ctx context.Context, req *authpb.ResetPasswordRequest) error {
	s.log.Info("Check if the user exists")
	formattedEmail := strings.ToLower(req.Email)
	user, err := s.userRepo.GetUserByEmail(ctx, formattedEmail)

	if err != nil {
		return status.Error(codes.NotFound, constants.ErrUserNotFound)
	}

	s.log.Info("Check if there is any active otp")
	otp, err := s.otpRepo.CheckActiveOtp(ctx, user.ID)
	if err != nil {
		s.log.Error("failed to get active otp", "error", err)
		return status.Error(codes.Unauthenticated, constants.ErrOtpNotFound)
	}

	if otp.Code != req.Code {
		s.log.Warn("invalid otp code")
		return status.Error(codes.Unauthenticated, constants.ErrOtpInvalid)
	}

	newPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		s.log.Error("failed to hash the password", "error", err)
		return status.Error(codes.Internal, constants.ErrInternalServer)
	}

	if err := s.userRepo.ResetPassword(ctx, user.ID, string(newPassword)); err != nil {
		s.log.Error("failed to reset the password", "error", err)
		return status.Error(codes.Internal, constants.ErrInternalServer)
	}

	if err := s.otpRepo.MarkOtpAsUsed(ctx, otp.ID); err != nil {
		s.log.Error("failed to reset the password", "error", err)
		return status.Error(codes.Internal, constants.ErrInternalServer)
	}

	event := map[string]any{
		"event":    "user.otp.passwordReset",
		"email":    user.Email,
		"fullName": user.FullName,
		"location": "Hargeisa, Somaliland", // FIXME: Update this to the req IP location
	}

	body, err := json.Marshal(event)
	if err != nil {
		s.log.Error("failed to parse otp event", "error", err)
		return status.Error(codes.Internal, constants.ErrOtpGenerationFail)
	}

	s.log.Info("Publish password reset event")
	if err := s.rabbit.Publish("user.otp.passwordReset", body); err != nil {
		s.log.Error("failed to publish password reset request event", "error", err)
		return status.Error(codes.Internal, constants.ErrOtpGenerationFail)
	}

	return nil

}
