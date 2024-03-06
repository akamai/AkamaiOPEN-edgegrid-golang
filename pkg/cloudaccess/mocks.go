//revive:disable:exported

package cloudaccess

import (
	"github.com/stretchr/testify/mock"
)

type Mock struct {
	mock.Mock
}

var _ CloudAccess = &Mock{}
