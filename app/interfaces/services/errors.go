package services

type HttpError interface {
	StatusCode() int
	String() string
}

type ApplicationError interface {
	StatusCode() int
	String() string
	StackTrace() string
}


