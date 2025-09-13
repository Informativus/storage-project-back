package user_service

import (
	"time"

	"github.com/google/uuid"
	"github.com/ivan/storage-project-back/internal/models/roles_model"
	"github.com/ivan/storage-project-back/internal/models/user_model"
	"github.com/ivan/storage-project-back/internal/repository/user_repo"
	"github.com/ivan/storage-project-back/internal/services/folder_service"
	"github.com/ivan/storage-project-back/internal/services/security_service"
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
	security      *security_service.SecurityService
}

func NewUserService(
	cfg *config.Config,
	jwt *jwt_service.JwtService,
	userRepo *user_repo.UserRepo,
	fldService *folder_service.FolderService,
	secService *security_service.SecurityService,
) *UserService {
	return &UserService{
		cfg:           cfg,
		jwt:           jwt,
		UserRepo:      userRepo,
		FolderService: fldService,
		security:      secService,
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
		RoleID:    roles_model.User,
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

func (u *UserService) MarkUsrAsDeleted(id uuid.UUID) error {
	err := u.UserRepo.DelUsrFromCache(id)

	if err != nil {
		return errsvc.UsrErr.Internal.New(err)
	}

	if tag, err := u.UserRepo.MarkUsrAsDeleted(id); err != nil || tag == 0 {
		return errsvc.UsrErr.Internal.New(err)
	}

	return nil
}

func (u *UserService) Me(id uuid.UUID) (*user_model.UserDto, error) {
	usrModelCache, err := u.UserRepo.GetUserByIdCache(id)

	if err != nil {
		return nil, errsvc.UsrErr.Internal.New(err)
	}

	if usrModelCache != nil {
		return usrModelCache, nil
	}

	usrModel, err := u.UserRepo.GetUserByIdDb(id)

	if err != nil {
		return nil, errsvc.UsrErr.Internal.New(err)
	}

	if usrModel == nil {
		return nil, errsvc.UsrErr.NotFound.New(err)
	}

	fldProtectionGroups, err := u.security.GetUserFldProtectionGroups(usrModel.ID)

	if err != nil {
		return nil, errsvc.UsrErr.Internal.New(err)
	}

	usrDto := &user_model.UserDto{
		ID:              usrModel.ID,
		Name:            usrModel.Name,
		Blocked:         usrModel.Blocked,
		RoleID:          usrModel.RoleID,
		FoldersToAccess: fldProtectionGroups,
		CreatedAt:       usrModel.CreatedAt,
		UpdatedAt:       usrModel.UpdatedAt,
		DeletedAt:       usrModel.DeletedAt,
	}

	err = u.UserRepo.SetUsrInCache(usrDto)

	return usrDto, err
}

func (u *UserService) UpdateBlockUserInf(usrName string, blocked bool) error {
	usrModel, err := u.UserRepo.GetUserByName(usrName)

	if err != nil {
		return errsvc.UsrErr.Internal.New(err)
	}

	if usrModel == nil {
		return errsvc.UsrErr.NotFound.New(err)
	}

	_, err = u.UserRepo.UpdateBlockUserInf(blocked, usrModel.ID)

	if err != nil {
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

func (u *UserService) GetUserAccessByToken(token string) (*user_model.UserTokensModel, error) {
	return u.UserRepo.GetUserAccessByToken(token)
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

	if delErr := u.deleteUsr(id); delErr != nil {
		log.Error().Err(delErr).Str("user_id", id.String()).
			Msg("rollback failed, inconsistent state: manual cleanup required")
		return errsvc.UsrErr.InconsistentState.New(err)
	}
	return err
}

func (u *UserService) deleteUsr(id uuid.UUID) error {
	if _, err := u.UserRepo.DelUser(id); err != nil {
		return err
	}

	return nil
}
