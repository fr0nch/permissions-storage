package permissions

// Generated from permissions

type Status = int32

const (
	Status_Success              Status = 0
	Status_Allow                Status = 1
	Status_Disallow             Status = 2
	Status_PermNotFound         Status = 3
	Status_CookieNotFound       Status = 4
	Status_GroupNotFound        Status = 5
	Status_ChildGroupNotFound   Status = 6
	Status_ParentGroupNotFound  Status = 7
	Status_ActorUserNotFound    Status = 8
	Status_TargetUserNotFound   Status = 9
	Status_GroupAlreadyExist    Status = 10
	Status_UserAlreadyExist     Status = 11
	Status_CallbackAlreadyExist Status = 12
	Status_CallbackNotFound     Status = 13
	Status_PermAlreadyGranted   Status = 14
	Status_TemporalGroup        Status = 15
	Status_PermanentGroup       Status = 16
	Status_GroupNotDefined      Status = 17
)

type Action = int32

const (
	Action_Add    Action = 0
	Action_Remove Action = 1
)
