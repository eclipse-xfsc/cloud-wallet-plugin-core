package core

import (
	"log"

	types "github.com/eclipse-xfsc/microservice-core-go/pkg/logr"
)

var logger = getLogger()

func getLogger() types.Logger {
	l, err := types.New(libConfig.LogLevel, libConfig.IsDev, nil)
	if err != nil {
		log.Fatal(err)
	}
	return *l
}
