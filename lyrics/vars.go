package lyrics

import "errors"

var ErrAPIDownError = errors.New("API service is not available at the moment\n..." +
	"try again later (report to support team if the problem insists)")

var ErrNameTooShort = errors.New("song name is too short...\n" +
	"please try to enter the song name with format <singer name> - <song name> ")

var ErrNotFound = errors.New("the requested query not found :(\n" +
	"please ensure that the song actualy exists")

var ErrGoodResultNotFound = errors.New("couldn't find a good result for this query...\n" +
	"please try to enter the song name with this format: <singer name> - <song name> ")

var ErrAPINotAvailable = errors.New("couldn't reach the API servers...\n" +
	"please report the problem to support team")
