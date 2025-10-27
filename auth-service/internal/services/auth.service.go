package services

import (
	"auth-service/internal/interceptor"
	"auth-service/internal/redisconn"
	db "auth-service/internal/repository"
	"auth-service/utils"
	"context"
	"database/sql"
	"errors"
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
	log.Println("Login request for username: ", req.Username)
	log.Println("password: ", req.Password)
	user, err := s.Queries.GetUserByUsername(ctx, req.Username)
	log.Println("User fetched: ", user)
	log.Println("2fa: ", user.Is2fa.Bool)

	if err != nil {
		log.Println("Error fetching user: ", err)
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}
	log.Println("Comparing password hash...")
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		log.Println("Password comparison failed: ", err)
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}
	log.Println("Password matched.")

	log.Println("Parsing device ID...")
	_, err = uuid.Parse(req.DeviceId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid device id")
	}
	log.Println("Device ID parsed: ", req.DeviceId)
	log.Println("Getting user agent...")
	userAgent, ok := utils.GetUserMetadata(ctx, interceptor.UserAgentKey)
	if !ok {
		log.Println("UserAgent not found")
	}
	browser, os := utils.ParseUA(userAgent)
	log.Println("Parsed user agent - Browser: ", browser, ", OS: ", os)
	var trusted int32
	trusted = 0
	log.Println("Checking 2fa...")
	if user.Is2fa.Bool {
		log.Println("2fa enabled")
		trusted, err = s.Queries.IsValidTrustedDevice(ctx, db.IsValidTrustedDeviceParams{
			UserID:   user.UserID,
			DeviceID: req.DeviceId,
			Browser:  browser,
			Os:       os,
		})
		if err != nil && errors.Is(err, sql.ErrNoRows) {
			log.Println("Error checking trusted device: ", err)
			return nil, status.Error(codes.Internal, "check trusted device failed")
		}
		log.Println("trusted: ", trusted)
		if trusted != 1 {
			if req.Verified_2FaSessionId == "" {
				log.Println("no verified session id, need 2fa")
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
	log.Println("2fa not enabled or device trusted or 2fa verified")
	accessJti := uuid.New()
	refreshJti := uuid.New()
	changePassAt := user.ChangePassAt.Time.Unix()
	go func() {
		if err := redisconn.SetChangePassAt(user.UserID, user.ChangePassAt.Time.Unix()); err != nil {
			log.Println("redis err: ", err)
		}
	}()

	accessToken, err := utils.GenerateAccessToken(user.UserID, accessJti, user.DisplayName.String, changePassAt)
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
	log.Println("login successful: ", user.UserID, " - trusted: ", trusted)
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
	log.Println("is ok: ", canUse, " - reason: ", reason)

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
	log.Println("VerifyCode called")
	log.Println("SessionId: ", req.SessionId)
	log.Println("DeviceId: ", req.DeviceId)
	log.Println("Otp: ", req.Otp)
	log.Println("RememberDevice: ", req.RememberDevice)
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
		IsSuccess: true,
		Message:   "OTP verified successfuly!",
	}, nil
}

func (s *AuthServer) Toggle2Fa(ctx context.Context, req *auth.Toggle2FaRequest) (*auth.Toggle2FaResponse, error) {
	log.Println("Toggle2Fa called")
	log.Println("Enable: ", req.Enable)
	log.Println("DeviceId: ", req.DeviceId)
	log.Println("Verified_2FaSessionId: ", req.Verified_2FaSessionId)
	userID, ok := utils.GetUserMetadata(ctx, interceptor.UserIDKey)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "user not authenticated")
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, status.Error(codes.Internal, "invalid user id")
	}
	userEmail := ""
	userEmailRaw, err := s.Queries.GetEmail(ctx, userUUID)
	if err != nil {
		return nil, status.Error(codes.Internal, "get user email failed")
	}

	if userEmailRaw.EmailEncrypted.String == "" {
		return nil, status.Error(codes.FailedPrecondition, "no email associated with account")
	}
	userEmail, err = utils.DecryptEmailRSA(userEmailRaw.EmailEncrypted.String)
	if err != nil {
		return nil, status.Error(codes.Internal, "decrypt email failed")
	}

	if req.Verified_2FaSessionId != "" {
		verified := redisconn.CheckVerifiedSession(req.Verified_2FaSessionId, req.DeviceId)
		if !verified {
			return nil, status.Error(codes.PermissionDenied, "2FA not verified")
		}
		err = s.Queries.Toggle2Fa(ctx, db.Toggle2FaParams{
			UserID: userUUID,
			Is2fa:  pgtype.Bool{Bool: req.Enable, Valid: true},
		})
		if err != nil {
			return nil, status.Error(codes.Internal, "toggle 2fa failed")
		}
		return &auth.Toggle2FaResponse{
			IsSuccess: true,
			Message:   "2FA toggled successfully",
		}, nil
	}

	verifing, err := redisconn.GetVerifing2Fa(userID)
	if err != nil && verifing == "1" {
		// logout all
		return nil, status.Error(codes.PermissionDenied, "Multiple security-sensitive actions detected. For your safety we've logged you out. Please re-login and verify via 2FA.")
	}
	err = redisconn.SetVerifing2Fa(userID)
	if err != nil {
		return nil, status.Error(codes.Internal, "set verifing 2fa failed")
	}
	sessionId := uuid.New()
	otp := utils.GenerateOTP(int(utils.TOGGLE_2FA_OTP))
	err = redisconn.Set2faOTP(sessionId.String(), userID, req.DeviceId, otp)
	if err != nil {
		return nil, status.Error(codes.Internal, "set 2fa otp failed")
	}

	err = utils.SendEmailOTP(userEmail, otp, "vi")
	if err != nil {
		log.Println("send email otp failed: ", err)
		return nil, status.Error(codes.Internal, "send email otp failed")
	}
	return &auth.Toggle2FaResponse{
		IsSuccess: true,
		SessionId: sessionId.String(),
		Message:   "OTP sent to your email",
	}, nil
}

func (s *AuthServer) GetProfile(ctx context.Context, req *auth.GetProfileRequest) (*auth.GetProfileResponse, error) {
	log.Println("GetProfile called")
	userId, ok := utils.GetUserMetadata(ctx, interceptor.UserIDKey)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "user not authenticated")
	}

	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return nil, status.Error(codes.Internal, "invalid user id")
	}

	user, err := s.Queries.GetUserByID(ctx, userUUID)
	if err != nil {
		return nil, status.Error(codes.Internal, "get user failed")
	}
	userEmail := ""
	if user.EmailEncrypted.String != "" {
		userEmail, err = utils.DecryptEmailRSA(user.EmailEncrypted.String)
		if err != nil {
			return nil, status.Error(codes.Internal, "decrypt email failed")
		}
	}

	return &auth.GetProfileResponse{
		AvatarUrl:   user.AvatarUrl.String,
		DisplayName: user.DisplayName.String,
		Is_2Fa:      user.Is2fa.Bool,
		Email:       userEmail,
		CreatedAt:   user.CreatedAt.Time.Unix(),
		UpdatedAt:   user.UpdatedAt.Time.Unix(),
	}, nil
}

func (s *AuthServer) UpdateProfile(ctx context.Context, req *auth.UpdateProfileRequest) (*auth.UpdateProfileResponse, error) {
	return nil, nil
}

func (s *AuthServer) CanUseEmail(ctx context.Context, email string) (bool, string) {
	emailHash := utils.HashEmail(email)
	userID, ok := utils.GetUserMetadata(ctx, interceptor.UserIDKey)
	if !ok {
		return false, "user not authenticated"
	}

	dbMail, err := s.Queries.CheckActiveEmail(ctx, pgtype.Text{String: emailHash, Valid: true})
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		log.Println("db error when checking active email: ", err)
		return false, "db error"
	}
	if dbMail.String != "" {
		return false, "email already in use"
	}

	deletedEmail, err := s.Queries.GetDeletedEmail(ctx, emailHash)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		log.Println("db error when checking deleted email: ", err)
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
