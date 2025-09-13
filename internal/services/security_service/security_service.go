package security_service

import (
	"github.com/google/uuid"
	"github.com/ivan/storage-project-back/internal/models/folder_model"
	"github.com/ivan/storage-project-back/internal/models/protection_group"
	"github.com/ivan/storage-project-back/internal/models/user_model"
	"github.com/ivan/storage-project-back/internal/repository/folder_repo"
	"github.com/ivan/storage-project-back/internal/repository/security_repo"
	"github.com/ivan/storage-project-back/pkg/errsvc"
)

type SecurityService struct {
	securityRepo *security_repo.SecurityRepo
	fldRepo      *folder_repo.FldRepo
}

func NewSecurityService(secRepo *security_repo.SecurityRepo, fldRepo *folder_repo.FldRepo) *SecurityService {
	return &SecurityService{
		securityRepo: secRepo,
		fldRepo:      fldRepo,
	}
}

func (ss *SecurityService) AccessToCreateFldForUsr(usrDto *user_model.UserDto, fldToCreateID uuid.UUID) error {
	fldAccessModel, err := ss.securityRepo.GetUsrRoleForFolder(usrDto.ID, fldToCreateID)

	if err != nil {
		return errsvc.SecurityErr.Internal.New(err)
	}

	if fldAccessModel == nil {
		return errsvc.SecurityErr.NotFound.New(err)
	}

	return nil
}

func (ss *SecurityService) AccessToUploadFileInFld(usrModel *user_model.UserModel, fldID uuid.UUID) error {
	fldAccessModel, err := ss.securityRepo.GetUsrRoleForFolder(usrModel.ID, fldID)

	if err != nil {
		return errsvc.SecurityErr.Internal.New(err)
	}

	if fldAccessModel == nil {
		return errsvc.SecurityErr.NotFound.New(err)
	}

	return nil
}

func (ss *SecurityService) AssignFolderGroups(fldAccModel *folder_model.FolderAccessModel, groups []protection_group.ProtectionGroupsType) error {
	fldAccModel, err := ss.securityRepo.InsertFolderAccess(fldAccModel)

	if err != nil {
		return errsvc.SecurityErr.Internal.New(err)
	}

	for _, group := range groups {
		groupID, ok := protection_group.ProtectionGroupIDs[group]

		if !ok {
			return errsvc.SecurityErr.Internal.New(err)
		}

		_, err := ss.securityRepo.InsertFolderProtectionGroups(&protection_group.FolderProtectionGroupsModel{
			FolderID: fldAccModel.FolderID,
			UserID:   fldAccModel.UserID,
			GroupID:  groupID,
		})

		if err != nil {
			return errsvc.SecurityErr.Internal.New(err)
		}
	}

	return nil
}

func (ss *SecurityService) GetUserFldProtectionGroups(usrID uuid.UUID) ([]folder_model.FolderAccessDto, error) {
	fldProtectGroups, err := ss.securityRepo.GetUsrFoldersProtectionGroups(usrID)

	if err != nil {
		return nil, errsvc.SecurityErr.Internal.New(err)
	}

	folderMap := make(map[uuid.UUID]*folder_model.FolderAccessDto)
	for _, fpg := range fldProtectGroups {
		dto, exist := folderMap[fpg.FolderID]

		if !exist {
			fld, err := ss.fldRepo.GetFldById(fpg.FolderID)

			if err != nil {
				return nil, errsvc.SecurityErr.Internal.New(err)
			}

			dto = &folder_model.FolderAccessDto{
				FolderID:         fpg.FolderID,
				FolderName:       fld.Name,
				ProtectionGroups: []protection_group.ProtectionGroupsType{},
			}
			folderMap[fpg.FolderID] = dto
		}

		dto.ProtectionGroups = append(dto.ProtectionGroups, protection_group.ProtectionGroupNames[fpg.GroupID])
	}

	var result []folder_model.FolderAccessDto
	for _, dto := range folderMap {
		result = append(result, *dto)
	}

	return result, nil
}

func (ss *SecurityService) HasAccess(usrDto *user_model.UserDto) {

}
