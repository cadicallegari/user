package xdatabase

import "github.com/cadicallegari/user/pkg/xerrors"

func NewConnectionError(driverName string, err error) error {
	if err == nil {
		return xerrors.Newf(xerrors.Internal, driverName+"_connection_error", "unable to connect to %s", driverName)
	}
	return xerrors.Newf(xerrors.Internal, driverName+"_connection_error", "unable to connect to %s: %w", driverName, err)
}
