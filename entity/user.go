package entity

type UserW struct {
	ID          string
	Username    string
	Password    string
	Nickname    string
	Bio         string
	AvatarPath  string
	CreatedAt   int64
	SuspendedAt int64
}

type UserR struct {
	ID          string
	Username    string
	Password    string
	Nickname    string
	Bio         string
	AvatarPath  string
	CreatedAt   int64
	SuspendedAt int64
}
