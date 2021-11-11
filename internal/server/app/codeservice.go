package app

import (
	"context"
	"github.com/hotafrika/griz-backend/internal/server/app/authtoken"
	"github.com/hotafrika/griz-backend/internal/server/app/password"
	"github.com/hotafrika/griz-backend/internal/server/app/token"
	"github.com/hotafrika/griz-backend/internal/server/domain"
	"github.com/hotafrika/griz-backend/internal/server/domain/entities"
	"github.com/hotafrika/griz-backend/internal/server/infrastructure/cache"
	"github.com/hotafrika/griz-backend/internal/server/infrastructure/instagram"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"strconv"
	"time"
)

// CodeService contains app logic
type CodeService struct {
	authTokenTTL       time.Duration
	hashTTL            time.Duration
	socialLinkTTL      time.Duration
	logger             zerolog.Logger
	cache              domain.Cacher
	codeRepo           domain.CodeRepository
	userRepo           domain.UserRepository
	qrSource           *instagram.QRSource
	passEncryptor      password.Encryptor
	authTokenEncryptor authtoken.JWT
	hashEncryptor      token.AES
}

// NewCodeService creates new service
func NewCodeService(
	authTokenTTL,
	hastTTL,
	socialLinkTTL time.Duration,
	logger zerolog.Logger,
	cache domain.Cacher,
	codeRepo domain.CodeRepository,
	userRepo domain.UserRepository,
	passEncryptor password.Encryptor,
	authTokenEncryptor authtoken.JWT,
	hashEncryptor token.AES,
) CodeService {
	return CodeService{
		authTokenTTL:       authTokenTTL,
		hashTTL:            hastTTL,
		socialLinkTTL:      socialLinkTTL,
		logger:             logger,
		cache:              cache,
		codeRepo:           codeRepo,
		userRepo:           userRepo,
		passEncryptor:      passEncryptor,
		authTokenEncryptor: authTokenEncryptor,
		hashEncryptor:      hashEncryptor,
		qrSource:           instagram.NewQRSource(),
	}
}

// CreateAuthToken returns userID by authToken if last exists
func (s CodeService) CreateAuthToken(ctx context.Context, user entities.User) (string, error) {
	encodedPass, err := s.passEncryptor.EncodeString(user.Password)
	if err != nil {
		return "", err
	}
	user.Password = string(encodedPass)

	id, err := s.userRepo.GetByUsernameAndPass(ctx, user)
	if err != nil {
		return "", err
	}

	authToken, err := s.authTokenEncryptor.MakeByID(id)
	if err != nil {
		return "", err
	}

	err = s.cache.Set(ctx, cache.AuthToken{Key: authToken}, strconv.FormatUint(id, 10), s.authTokenTTL)
	if err != nil {
		return "", err
	}

	return authToken, nil
}

// GetUserIDByAuthToken returns userID by authToken if last exists
func (s CodeService) GetUserIDByAuthToken(ctx context.Context, authToken string) (uint64, error) {
	res, err := s.cache.Get(ctx, cache.AuthToken{Key: authToken})
	if err != nil {
		return 0, err
	}
	userID, err := strconv.ParseUint(res, 10, 64)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

// GetUser returns user by userID
func (s CodeService) GetUser(ctx context.Context, id uint64) (entities.User, error) {
	user, err := s.userRepo.Get(ctx, id)
	if err == nil {
		user.Password = ""
	}
	return user, err
}

// FindCodeBySocial returns sourceUrl by social link
func (s CodeService) FindCodeBySocial(ctx context.Context, link string) (string, error) {
	value, err := s.cache.Get(ctx, cache.SocialUrl{Key: link})
	if err == nil { // link found
		return value, nil
	}
	if !errors.Is(err, domain.ErrCacheNotExist) { // some error
		// TODO wrap
		return "", err
	}

	// link not found
	b, err := s.qrSource.GetFirstQR(ctx, link)
	if err != nil {
		// TODO wrap
		return "", err
	}

	hashToken, err := token.ExtractHashFromLink(string(b))
	if err != nil { // token not found in decoded data
		return "", errors.Wrap(err, "hash from decoded data: ")
	}

	srcLink, err := s.FindCodeByHash(ctx, hashToken)
	if err != nil {
		return "", errors.Wrap(err, "hash from decoded data: ")
	}

	err = s.cache.Set(ctx, cache.SocialUrl{Key: link}, srcLink, s.socialLinkTTL)
	if err != nil {
		// TODO wrap or maybe log
		return "", err
	}

	return string(b), nil
}

// FindCodeByHash returns sourceUrl by its hash
func (s CodeService) FindCodeByHash(ctx context.Context, hashToken string) (string, error) {
	value, err := s.cache.Get(ctx, cache.HashUrl{Key: hashToken})
	if err == nil { // hashToken found
		return value, nil
	}
	if !errors.Is(err, domain.ErrCacheNotExist) { // some error
		// TODO wrap
		return "", err
	}

	// hashToken not found
	code, err := s.codeRepo.GetByHash(ctx, hashToken)
	if err != nil {
		// TODO wrap
		return "", err
	}

	err = s.cache.Set(ctx, cache.HashUrl{Key: code.Hash}, code.SrcURL, s.hashTTL)
	if err != nil {
		// TODO wrap or log
		return "", err
	}
	return code.SrcURL, nil
}

// CreateCode creates code and adds it to cache
func (s CodeService) CreateCode(ctx context.Context, code entities.Code) (uint64, error) {
	hashValue, err := s.hashEncryptor.Create(code.ID)
	if err != nil {
		return 0, err
	}

	code.Hash = hashValue
	id, err := s.codeRepo.Create(ctx, code)
	if err != nil {
		return 0, err
	}

	err = s.cache.Set(ctx, cache.HashUrl{Key: hashValue}, strconv.FormatUint(id, 10), s.hashTTL)
	if err != nil {
		// TODO maybe delete from repo
		return 0, err
	}

	return id, nil
}
