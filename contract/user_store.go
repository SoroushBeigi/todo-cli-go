package contract

import "github.com/SoroushBeigi/todo-cli-go/entity"

type UserWriteStore interface {
	Save(u entity.User)
}

type UserReadStore interface {
	Load() []entity.User
}
