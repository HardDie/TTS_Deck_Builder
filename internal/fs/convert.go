package fs

import (
	"strconv"

	"github.com/HardDie/DeckBuilder/internal/errors"
)

func StringToInt64(in string) (int64, error) {
	val, err := strconv.ParseInt(in, 10, 64)
	if err != nil {
		errors.IfErrorLog(err)
		return 0, errors.BadId.AddMessage(err.Error())
	}
	return val, nil
}
