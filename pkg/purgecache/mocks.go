//revive:disable:exported

package purgecache

import "github.com/stretchr/testify/mock"

type Mock struct {
	mock.Mock
}

var _ PurgeCache = &Mock{}
