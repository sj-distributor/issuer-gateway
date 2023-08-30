package utils

import (
	"github.com/yitter/idgenerator-go/idgen"
)

// Id  Generate an ID.
func Id() int64 {
	return idgen.NextId()
}

func init() {
	idgen.SetIdGenerator(idgen.NewIdGeneratorOptions(1))
}
