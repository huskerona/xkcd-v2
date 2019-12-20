package model

type ValidationStatus struct {
	Id           int
	MissingJSON  bool
	MissingImage bool
}
