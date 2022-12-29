package translationMapping

import (
	"github.com/CanadianCommander/translationBot/internal/lib/translation"
	"io"
)

//==========================================================================
// Public
//==========================================================================

type MappingLoader interface {

	// Load - load the given file in to the standard format.
	Load(fileData io.Reader) ([]translation.Translation, error)

	// IsMimeSupported - check if the given mime type is supported by this loader
	// #### params
	// mimetype - the mime type to check
	IsMimeSupported(mimetype string) bool

	// IsFileSupported - check if the given file data is supported by this loader. When this is invoked the mime type
	// check will have already succeeded.
	// #### params
	// fileDataStart - the file data to check for compatability. Contains, at max, the first 128KB of the file
	IsFileSupported(fileDataStart []byte) bool

	// GetLoaderType returns a loader type string.
	GetLoaderType() string
}
