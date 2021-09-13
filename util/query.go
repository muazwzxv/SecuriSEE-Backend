package util

type Query struct {
	Page  string `query:"page"`
	Skip  string `query:"skip"`
	Limit string `query:"limit"`
}
