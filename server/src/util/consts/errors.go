package consts

import "errors"

var (
	ErrTwitterSearchParamPagesTooBig   = errors.New("params 'page' is too big for twitter search api")
	ErrTwitterSearchParamPagesTooSmall = errors.New("params 'page' is too small for twitter search api")
)
