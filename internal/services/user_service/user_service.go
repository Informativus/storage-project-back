package user_service

import (
	"time"

	"github.com/google/uuid"
	"github.com/ivan/storage-project-back/internal/models/roles_model"
	"github.com/ivan/storage-project-back/internal/models/user_model"
	"github.com/ivan/storage-project-back/internal/repository/user_repo"
	"github.com/ivan/storage-project-back/internal/services/folder_service"
	"github.com/ivan/storage-project-back/pkg/config"
	"github.com/ivan/storage-project-back/pkg/errsvc"
	"github.com/ivan/storage-project-back/pkg/jwt_service"
	"github.com/rs/zerolog/log"
)

type UserService struct {
	cfg           *config.Config
	jwt           *jwt_service.JwtService
	UserRepo      *user_repo.UserRepo
	FolderService *folder_service.FolderService
}

func NewUserService(
	cfg *config.Config,
	jwt *jwt_service.JwtService,
	userRepo *user_repo.UserRepo,
	fldService *folder_service.FolderService,
) *UserService {
	return &UserService{
		cfg:           cfg,
		jwt:           jwt,
		UserRepo:      userRepo,
		FolderService: fldService,
	}
}

func (u *UserService) CreateUser(usrName string) (string, error) {
	folderExt, err := u.FolderService.MainFolderExist(usrName)

	if err != nil {
		return "", errsvc.UsrErr.Internal.New(err)
	}

	if folderExt {
		return "", errsvc.UsrErr.AlreadyExists.New(err)
	}

	usrID := uuid.New()

	usrModel, err := u.UserRepo.CreateUser(user_model.UserModel{
		ID:        usrID,
		Name:      usrName,
		Blocked:   false,
		RoleID:    roles_model.Admin,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if err != nil {
		return "", errsvc.UsrErr.Internal.New(err)
	}

	token, err := u.addUserToken(usrModel)

	if err != nil {
		return "", u.rollbackUser(usrModel.ID, "failed to create user access token", err)
	}

	err = u.FolderService.CreateFolder(usrName, usrModel)

	if err != nil {
		return "", u.rollbackUser(usrModel.ID, "failed to create user folder", err)
	}

	return token, nil
}

func (u *UserService) DelUser(id uuid.UUID) error {
	usr, err := u.UserRepo.GetUserById(id)

	if err != nil {
		return errsvc.UsrErr.BadReq.New(err)
	}

	if usr == nil {
		return errsvc.UsrErr.NotFound.New(err)
	}

	if tag, err := u.UserRepo.DelUser(id); err != nil || tag == 0 {
		return errsvc.UsrErr.Internal.New(err)
	}

	return nil
}

func (u *UserService) AddUserTokenByUsrName(usrName string) (string, error) {
	usrModel, err := u.UserRepo.GetUserByName(usrName)

	if err != nil {
		return "", errsvc.UsrErr.Internal.New(err)
	}

	if usrModel == nil {
		return "", errsvc.UsrErr.NotFound.New(err)
	}

	return u.addUserToken(usrModel)
}

func (u *UserService) addUserToken(usrModel *user_model.UserModel) (string, error) {
	token, err := u.generateToken(jwt_service.JwtPayload{
		ID: usrModel.ID,
	})

	if err != nil {
		return "", err
	}

	usrAccessModel, err := u.UserRepo.InsertUserToken(user_model.UserTokensModel{
		ID:        uuid.New(),
		UserID:    usrModel.ID,
		Token:     token,
		Revoked:   false,
		CreatedAt: usrModel.CreatedAt,
		ExpiresAt: usrModel.CreatedAt.Add(time.Duration(u.cfg.ExpiresIn) * time.Second),
	})

	if err != nil {
		return "", errsvc.UsrErr.Internal.New(err)
	}

	return usrAccessModel.Token, nil
}

func (u *UserService) generateToken(payload jwt_service.JwtPayload) (string, error) {
	token, err := u.jwt.GenerateToken(payload)

	if err != nil {
		return "", errsvc.UsrErr.GenerateToken.New(err)
	}

	return token, nil
}

func (u *UserService) rollbackUser(id uuid.UUID, logMsg string, err error) error {
	log.Error().Err(err).Msg(logMsg)

	if delErr := u.DelUser(id); delErr != nil {
		log.Error().Err(delErr).Str("user_id", id.String()).
			Msg("rollback failed, inconsistent state: manual cleanup required")
		return errsvc.UsrErr.InconsistentState.New(err)
	}
	return err
}
