package logging

import (
	"github.com/ALiwoto/ssg/ssg"
	"go.uber.org/zap"
)

var AppLogger *zap.SugaredLogger

// LogStorages is a list of all log storages.
// by default, there is a single log storage which
// uses memory. The database package can add another
// log storage to this array.
var LogStorages []LogStorage

var logIdGenerator = ssg.NewNumIdGenerator[int64]()
