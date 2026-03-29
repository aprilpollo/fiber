package models

type TableNamer interface {
	TableName() string
}

type ModelList []TableNamer

func All() ModelList {
	return ModelList{
		&UserModel{},
		&BookModel{},
	}
}