//revive:disable:exported

package cloudwrapper

import (
	"github.com/stretchr/testify/mock"
)

type Mock struct {
	mock.Mock
}

var _ CloudWrapper = &Mock{}
