package app

import (
	"context"
	"encoding/base64"
	"github.com/hotafrika/griz-backend/internal/server/app/authtoken"
	"github.com/hotafrika/griz-backend/internal/server/app/password"
	"github.com/hotafrika/griz-backend/internal/server/app/qrencoder"
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
	qrEncoder          qrencoder.Yeqown
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
		qrEncoder:          qrencoder.DefaultYeqown(),
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
		return 0, errors.Wrap(err, "GetUserIDByAuthToken: ")
	}
	userID, err := strconv.ParseUint(res, 10, 64)
	if err != nil {
		return 0, errors.Wrap(err, "GetUserIDByAuthToken: ")
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
	id, err := s.codeRepo.Create(ctx, code)
	if err != nil {
		return 0, err
	}

	hashValue, err := s.hashEncryptor.Create(id)
	if err != nil {
		return 0, err
	}

	code.Hash = hashValue
	code.ID = id

	err = s.codeRepo.Update(ctx, code)
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

// GetCodes returns codes by userID
func (s CodeService) GetCodes(ctx context.Context, userID uint64) ([]entities.Code, error) {
	codes, err := s.codeRepo.ListAll(ctx, userID)
	return codes, err
}

// GetCode returns code by userID
func (s CodeService) GetCode(ctx context.Context, codeID uint64) (entities.Code, error) {
	code, err := s.codeRepo.Get(ctx, codeID)
	return code, err
}

// DownloadCodeByHash returns code by userID
func (s CodeService) DownloadCodeByHash(ctx context.Context, hashToken string) (string, error) {
	b, err := s.qrEncoder.Encode([]byte(token.BuildLink(hashToken)))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

// UpdateCode ...
func (s CodeService) UpdateCode(ctx context.Context, code entities.Code) error {
	err := s.cache.Set(ctx, cache.HashUrl{Key: code.Hash}, code.SrcURL, s.hashTTL)
	if err != nil {
		return err
	}

	err = s.codeRepo.Update(ctx, code)
	return err
}

// DeleteCode ...
func (s CodeService) DeleteCode(ctx context.Context, code entities.Code) error {
	err := s.cache.Delete(ctx, cache.HashUrl{Key: code.Hash})
	if err != nil {
		return err
	}

	err = s.codeRepo.Delete(ctx, code.ID)
	return err
}
