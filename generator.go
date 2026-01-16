package jsonapi

type IDGenerator interface {
	Generate() (string, error)
}

const generateError = "failed to generate new identifier"

const notNilGeneratorError = "generator must be not nil"
