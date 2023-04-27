package handlers

import (
	"auth/domain/command"
	"auth/infrastructure/worker"
	"auth/service"
	"reflect"
)

var COMMAND_HANDLERS = map[reflect.Type]func(uow *service.UnitOfWork, mailer worker.WorkerInterface, cmd interface{}) error{
	reflect.TypeOf(&command.CreateTeam{}):             CreateTeamWrapper,
	reflect.TypeOf(&command.UpdateTeam{}):             UpdateTeamWrapper,
	reflect.TypeOf(&command.UpdateLastActiveTeam{}):   UpdateLastActiveTeamWrapper,
	reflect.TypeOf(&command.DeleteTeamMember{}):       DeleteTeamMemberWrapper,
	reflect.TypeOf(&command.ChangeMemberRole{}):       ChangeMemberRoleWrapper,
	reflect.TypeOf(&command.DeleteUser{}):             DeleteUserWrapper,
	reflect.TypeOf(&command.UpdateUser{}):             UpdateUserWrapper,
	reflect.TypeOf(&command.InviteMember{}):           InviteMemberWrapper,
	reflect.TypeOf(&command.ResendInvitation{}):       ResendInvitationWrapper,
	reflect.TypeOf(&command.DeleteInvitation{}):       DeleteInvitationWrapper,
	reflect.TypeOf(&command.UpdateInvitationStatus{}): UpdateInvitationStatusWrapper,
	reflect.TypeOf(&command.LoginByGoogle{}):          LoginByGoogleWrapper,
	reflect.TypeOf(&command.UpdateTeamAvatar{}):       UpdateTeamAvatarWrapper,
	reflect.TypeOf(&command.DeleteTeamAvatar{}):       DeleteTeamAvatarWrapper,
	reflect.TypeOf(&command.UpdateUserAvatar{}):       UpdateUserAvatarWrapper,
	reflect.TypeOf(&command.DeleteUserAvatar{}):       DeleteUserAvatarWrapper,
}
