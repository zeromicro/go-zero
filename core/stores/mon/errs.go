package mon

import (
	"errors"

	"github.com/zeromicro/go-zero/core/errorx"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/x/mongo/driver/session"
)

func acceptable(err error) bool {
	return err == nil || isDupKeyError(err) ||
		errorx.In(err, mongo.ErrNoDocuments, mongo.ErrNilValue,
			mongo.ErrNilDocument, mongo.ErrNilCursor, mongo.ErrEmptySlice,
			// session errors
			session.ErrSessionEnded, session.ErrNoTransactStarted, session.ErrTransactInProgress,
			session.ErrAbortAfterCommit, session.ErrAbortTwice, session.ErrCommitAfterAbort,
			session.ErrUnackWCUnsupported, session.ErrSnapshotTransaction)
}

func isDupKeyError(err error) bool {
	var e mongo.WriteException
	if !errors.As(err, &e) {
		return false
	}

	return e.HasErrorCode(duplicateKeyCode)
}
