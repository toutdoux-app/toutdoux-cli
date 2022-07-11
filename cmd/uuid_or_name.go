package cmd

import (
	"github.com/gofrs/uuid"
)

const (
	uuidOrNameIsUUID = iota
	uuidOrNameIsName
)

type uuidOrName struct {
	value string
	uuid  uuid.UUID
}

func newUUIDOrName(value string) uuidOrName {
	un := uuidOrName{
		value: value,
	}

	un.tryUUID()

	return un
}

func (u uuidOrName) Type() int {
	if u.uuid != uuid.Nil {
		return uuidOrNameIsUUID
	}
	return uuidOrNameIsName
}

func (u *uuidOrName) tryUUID() {
	uuidV, err := uuid.FromString(u.value)
	if err != nil {
		return
	}

	u.uuid = uuidV
}
