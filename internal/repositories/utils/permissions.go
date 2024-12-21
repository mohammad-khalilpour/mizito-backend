package utils




type PermissionHandler interface {
	CheckUserHasAccessToTeam(userID uint, teamID uint) bool
	CheckUserHasAccessToProject(userID uint, projectID uint) bool
	CheckUserHasAccessToTask(userID uint, TaskID uint) bool
}


func NewPermissionHandler() PermissionHandler{
	return &permissionHandler{}
}


type permissionHandler struct {

}


func (ph *permissionHandler) CheckUserHasAccessToTeam(userID uint, teamID uint) bool{
	return false
}


func (ph *permissionHandler) CheckUserHasAccessToProject(userID uint, projectID uint) bool{
	return false
}


func (ph *permissionHandler) CheckUserHasAccessToTask(userID uint, TaskID uint) bool{
	return false
}