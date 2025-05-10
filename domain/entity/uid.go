package entity

import "github.com/google/uuid"

type UID struct {
	value string
}

func NewUID() UID {
	// TODO: ドメインがID生成の詳細を知っている状態なので抽出したい
	uuid := uuid.New()
	return UID{
		value: uuid.String(),
	}
}

func (u UID) String() string {
	return u.value
}

func (u UID) IsEqual(target UID) bool {
	return u.value == target.value
}
