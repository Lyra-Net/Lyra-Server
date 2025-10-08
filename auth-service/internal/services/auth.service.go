package services

import (
	"auth-service/internal/interceptor"
	"auth-service/internal/redisconn"
	db "auth-service/internal/repository"
	"auth-service/utils"
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/trandinh0506/BypassBeats/proto/gen/auth"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const TRUSTED_DEVICE_TIME_EXPIRE_TIME time.Duration = 30 * 24 * time.Hour

type AuthServer struct {
	auth.UnimplementedAuthServiceServer
	Queries *db.Queries
}

func NewAuthServer(q *db.Queries) *AuthServer {
	return &AuthServer{Queries: q}
}

func (s *AuthServer) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), 14)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "password hashing failed")
	}

	userID := uuid.New()
	err = s.Queries.CreateUser(ctx, db.CreateUserParams{
		UserID:       userID,
		Username:     req.Username,
		DisplayName:  pgtype.Text{String: req.DisplayName, Valid: req.DisplayName != ""},
		PasswordHash: string(hashed),
	})
	if err != nil {
		return nil, status.Errorf(codes.AlreadyExists, "username already in use")
	}
	// --- Kafka emit section ---
	// Extract user metadata
	// browser, OS := utils.ParseUA(GetUserMetadata(ctx, interceptor.UserAgentKey))
	// userIP := utils.ParseIP(GetUserMetadata(ctx, interceptor.UserIpKey))

	// Emit event to Kafka (security topic)
	// event := dto.UserRegistered{
	// 	UserID:    userID.String(),
	// 	DeviceID:  req.DeviceId,
	// 	Browser:   browser,
	// 	OS:        OS,
	// 	UserIP:    userIP,
	// 	Timestamp: time.Now().Unix(),
	// }

	// Non-blocking emit
	// go func() {
	// 	if err := s.producer.Emit(ctx, "security.user_registered", event); err != nil {
	// 		s.logger.Warn("failed to emit security register event", "error", err)
	// 	}
	// }()
	// --------------------------
	return &auth.RegisterResponse{Message: "User created"}, nil
}

func (s *AuthServer) Login(ctx context.Context, req *auth.LoginRequest) (*auth.AuthResponse, error) {
	user, err := s.Queries.GetUserByUsername(ctx, req.Username)

	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}

	_, err = uuid.Parse(req.DeviceId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid device id")
	}

	userAgent, ok := utils.GetUserMetadata(ctx, interceptor.UserAgentKey)
	if !ok {
		log.Println("UserAgent not found")
	}

	browser, os := utils.ParseUA(userAgent)

	var trusted int32
	trusted = 0

	if user.Is2fa.Bool {
		trusted, err = s.Queries.IsValidTrustedDevice(ctx, db.IsValidTrustedDeviceParams{
			UserID:   user.UserID,
			DeviceID: req.DeviceId,
			Browser:  browser,
			Os:       os,
		})
		if err != nil && err != sql.ErrNoRows {
			return nil, status.Error(codes.Internal, "check trusted device failed")
		}

		if trusted != 1 {
			if req.Verified_2FaSessionId == "" {
				verifing, err := redisconn.GetVerifing2Fa(user.UserID.String())
				if err != nil && verifing == "1" {
					// logout all
					return nil, status.Error(codes.PermissionDenied, "Multiple security-sensitive actions detected. For your safety we've logged you out. Please re-login and verify via 2FA.")
				}
				err = redisconn.SetVerifing2Fa(user.UserID.String())
				if err != nil {
					return nil, status.Error(codes.Internal, "set verifing 2fa failed")
				}
				sessionId := uuid.New()
				userEmail, err := utils.DecryptEmailRSA(user.EmailEncrypted.String)
				if err != nil {
					return nil, status.Error(codes.Internal, "decrypt email failed")
				}
				otp := utils.GenerateOTP(int(utils.TWO_FA_LOGIN_OTP))
				redisconn.Set2faOTP(sessionId.String(), user.UserID.String(), req.DeviceId, otp)
				utils.SendEmailOTP(userEmail, otp, "vi")
				return nil, status.Errorf(codes.PermissionDenied, "2FA_REQUIRED:SESSION:%s", sessionId)
			}

			ok := redisconn.CheckVerifiedSession(req.Verified_2FaSessionId, req.DeviceId)
			if !ok {
				return nil, status.Error(codes.PermissionDenied, "2FA not verified")
			}
		}
	}

	accessJti := uuid.New()
	refreshJti := uuid.New()
	changePassAt := user.ChangePassAt.Time.Unix()
	go func() {
		if err := redisconn.SetChangePassAt(user.UserID, user.ChangePassAt.Time.Unix()); err != nil {
			log.Println("redis err: ", err)
		}
	}()

	accessToken, err := utils.GenerateAccessToken(user.UserID, accessJti, changePassAt)
	if err != nil {
		return nil, status.Error(codes.Internal, "generate access token failed")
	}
	refreshToken, err := utils.GenerateRefreshToken(user.UserID, refreshJti)
	if err != nil {
		return nil, status.Error(codes.Internal, "generate refresh token failed")
	}

	err = s.Queries.CreateRefreshToken(ctx, db.CreateRefreshTokenParams{
		ID:        refreshJti,
		UserID:    user.UserID,
		Token:     refreshToken,
		AccessJti: accessJti,
		DeviceID:  pgtype.Text{String: req.DeviceId, Valid: req.DeviceId != ""},
		Browser:   pgtype.Text{String: browser, Valid: browser != ""},
		Os:        pgtype.Text{String: os, Valid: os != ""},
		ExpiresAt: pgtype.Timestamp{Time: time.Now().Add(7 * 24 * time.Hour), Valid: true},
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "store refresh token failed")
	}
	// --- Kafka emit section ---
	// userIp, ok := utils.GetUserMetadata(ctx, interceptor.UserIpKey)
	// if !ok {
	// 	log.Println("UserIp not found")
	// }
	// event := dto.UserLogin{
	// 	UserID:    userID.String(),
	// 	DeviceID:  req.DeviceId,
	// 	Browser:   browser,
	// 	OS:        os,
	// 	UserIP:    userIp,
	// 	Timestamp: time.Now().Unix(),
	// }
	// Non-blocking emit
	// go func() {
	// 	if err := s.producer.Emit(ctx, "security.user_login", event); err != nil {
	// 		s.logger.Warn("failed to emit security login event", "error", err)
	// 	}
	// }()
	// --------------------------
	//
	return &auth.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthServer) Logout(ctx context.Context, req *auth.LogoutRequest) (*auth.LogoutResponse, error) {
	// claims, err := utils.ParseRefreshToken(req.RefreshToken)
	// if err != nil{
	// 	return nil, status.Error(codes.Unauthenticated, "invalid refresh token")
	// }

	return nil, nil
}

func (s *AuthServer) RefreshToken(ctx context.Context, req *auth.RefreshTokenRequest) (*auth.AuthResponse, error) {
	return nil, nil
}

func (s *AuthServer) ChangePassword(ctx context.Context, req *auth.ChangePasswordRequest) (*auth.ChangePasswordResponse, error) {
	return nil, nil
}

func (s *AuthServer) ForgotPassword(ctx context.Context, req *auth.ForgotPasswordRequest) (*auth.ForgotPasswordResponse, error) {
	return nil, nil
}

func (s *AuthServer) AddEmail(ctx context.Context, req *auth.AddEmailRequest) (*auth.AddEmailResponse, error) {

	userID, ok := utils.GetUserMetadata(ctx, interceptor.UserIDKey)

	if !ok {
		return nil, status.Error(codes.Unauthenticated, "user not authenticated")
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, status.Error(codes.Internal, "invalid user id")
	}

	if req.Verified_2FaSessionId != "" {
		verified := redisconn.CheckVerifiedSession(req.Verified_2FaSessionId, req.DeviceId)
		if !verified {
			return nil, status.Error(codes.PermissionDenied, "2FA not verified")
		}
		emailEncrypted, err := utils.EncryptEmailRSA(req.Email)
		emailHash := utils.HashEmail(req.Email)
		if err != nil {
			return nil, status.Error(codes.Internal, "encrypt email failed")
		}

		s.Queries.AddEmail(ctx, db.AddEmailParams{
			UserID:         userUUID,
			EmailEncrypted: pgtype.Text{String: emailEncrypted, Valid: emailEncrypted != ""},
			EmailHash:      pgtype.Text{String: emailHash, Valid: emailHash != ""},
		})
		return &auth.AddEmailResponse{
			IsSuccess: true,
			SessionId: "",
		}, nil
	}

	canUse, reason := s.CanUseEmail(ctx, req.Email)
	log.Println("is ok: ", ok, " - reason: ", reason)

	if !canUse {
		return nil, status.Error(codes.InvalidArgument, reason)
	}
	if err := redisconn.SetVerifing2Fa(userID); err != nil {
		return nil, status.Error(codes.Internal, "can not set verifing session")
	}

	sessionId := uuid.New()
	otp := utils.GenerateOTP(int(utils.ADD_EMAIL_OTP))
	utils.SendEmailOTP(req.Email, otp, "vi")
	if err := redisconn.Set2faOTP(sessionId.String(), userID, req.DeviceId, otp); err != nil {
		return nil, status.Error(codes.Internal, "can not set 2fa otp")
	}

	return &auth.AddEmailResponse{
		IsSuccess: true,
		SessionId: sessionId.String(),
	}, nil
}

func (s *AuthServer) RemoveEmail(ctx context.Context, req *auth.RemoveEmailRequest) (*auth.RemoveEmailResponse, error) {
	return nil, nil
}

func (s *AuthServer) ResendVerification(ctx context.Context, req *auth.ResendVerificationRequest) (*auth.ResendVerificationResponse, error) {
	return nil, nil
}

func (s *AuthServer) VerifyCode(ctx context.Context, req *auth.VerifyCodeRequest) (*auth.VerifyCodeResponse, error) {
	storedOtp, userID, err := redisconn.Get2faOTP(req.SessionId, req.DeviceId)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "wrong OTP or expired")
	}
	if storedOtp != req.Otp {
		return nil, status.Error(codes.Unauthenticated, "wrong OTP or expired")
	}
	err = redisconn.SetVerifiedSession(req.SessionId, req.DeviceId)
	if err != nil {
		return nil, status.Error(codes.Internal, "can not set verified session")
	}
	err = redisconn.RemoveVerifing2Fa(userID)
	if err != nil {
		log.Println("can not remove verifing session, err: ", err)
		//
	}

	if req.RememberDevice {
		ua, ok := utils.GetUserMetadata(ctx, interceptor.UserAgentKey)
		if !ok {
			log.Println("can not get user agent")
		}
		browser, os := utils.ParseUA(ua)
		s.Queries.CreateOrUpdateTrustedDevice(ctx, db.CreateOrUpdateTrustedDeviceParams{
			UserID:    uuid.MustParse(userID),
			DeviceID:  req.DeviceId,
			Browser:   browser,
			Os:        os,
			ExpiresAt: pgtype.Timestamp{Time: time.Now().Add(TRUSTED_DEVICE_TIME_EXPIRE_TIME), Valid: true},
		})
	}
	return &auth.VerifyCodeResponse{
		Success: true,
		Message: "OTP verified successfuly!",
	}, nil
}

func (s *AuthServer) CanUseEmail(ctx context.Context, email string) (bool, string) {
	emailHash := utils.HashEmail(email)
	userID, ok := utils.GetUserMetadata(ctx, interceptor.UserIDKey)
	if !ok {
		return false, "user not authenticated"
	}

	dbMail, err := s.Queries.CheckActiveEmail(ctx, pgtype.Text{String: emailHash, Valid: true})
	if err != nil && err != sql.ErrNoRows {
		log.Println("db error: ", err)
		return false, "db error"
	}
	if dbMail.String != "" {
		return false, "email already in use"
	}

	deletedEmail, err := s.Queries.GetDeletedEmail(ctx, emailHash)
	if err != nil && err != sql.ErrNoRows {
		log.Println("db error: ", err)
		return false, "db error"
	}
	if deletedEmail.EmailHash != "" {
		now := time.Now()
		if deletedEmail.DeletedBy == uuid.MustParse(userID) {
			if deletedEmail.CooldownUntil.Time.After(now) {
				return false, fmt.Sprintf("you deleted this email, wait until %s", deletedEmail.CooldownUntil.Time.Format("YYYY-MM-DD HH:mm:ss"))
			}
		} else {
			if deletedEmail.SafeWindowUntil.Time.After(now) {
				return false, fmt.Sprintf("email was recently removed by another user, wait until %s", deletedEmail.SafeWindowUntil.Time.Format("YYYY-MM-DD HH:mm:ss"))
			}
		}
	}

	return true, ""
}
