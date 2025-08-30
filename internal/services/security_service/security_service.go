package security_service

import (
	"github.com/google/uuid"
	"github.com/ivan/storage-project-back/internal/models/folder_model"
	"github.com/ivan/storage-project-back/internal/models/roles_model"
	"github.com/ivan/storage-project-back/internal/models/user_model"
	"github.com/ivan/storage-project-back/internal/repository/security_repo"
	"github.com/ivan/storage-project-back/pkg/errsvc"
)

type SecurityService struct {
	SecurityRepo *security_repo.SecurityRepo
}

func NewSecurityService(repo *security_repo.SecurityRepo) *SecurityService {
	return &SecurityService{
		SecurityRepo: repo,
	}
}

func (ss *SecurityService) AccessToCreateFldForUsr(usrModel *user_model.UserModel, fldToCreateID uuid.UUID) error {
	fldAccessModel, err := ss.SecurityRepo.GetUsrRoleForFolder(usrModel.ID, fldToCreateID)

	if err != nil {
		return errsvc.SecurityErr.Internal.New()
	}

	if fldAccessModel == nil {
		return errsvc.SecurityErr.NotFound.New()
	}

	accessRoles := map[roles_model.Role]struct{}{
		roles_model.Owner: {},
		roles_model.Admin: {},
	}

	if _, ok := accessRoles[fldAccessModel.RoleID]; !ok {
		return errsvc.UsrErr.Forbidden.New()
	}

	return nil
}

func (ss *SecurityService) GetUsrRoleForFld(usrModel *user_model.UserModel, mainFld *folder_model.FolderModel) roles_model.Role {
	if mainFld.OwnerID == usrModel.ID {
		return roles_model.Owner
	}

	return roles_model.Editor
}
