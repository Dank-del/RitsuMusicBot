package lyrics

// url-related constants
const (
	baseUrl    = "https://api.genius.com/"
	searchMeth = "search?q="
)

// character region constants
const (
	startChars1     = 8208
	endChars1       = 8231
	endEnglishChars = 126

	weirdChar1 = '\u00a0'
	weirdChar2 = '\u200b'
)

// hit-related constant region :3
const (
	minConfidence       = 2
	differentLangsPoint = 5

	songValue         = "song"
	completeValue     = "complete"
	instrumentalValue = "instrumental"
)

// unimportant string constants
const (
	unimportant1 = "and"
	unimportant2 = "or"
	unimportant3 = "to"
	unimportant4 = "so"
)
