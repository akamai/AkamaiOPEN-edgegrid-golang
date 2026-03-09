//revive:disable:exported

package reportinggroups

import "github.com/stretchr/testify/mock"

type Mock struct {
	mock.Mock
}

var _ ReportingGroups = &Mock{}
