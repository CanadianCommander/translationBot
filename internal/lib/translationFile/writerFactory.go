package translationFile

//==========================================================================
// Public
//==========================================================================

// GetWriterForFile returns the appropriate translation writer for the specified file
// #### params
// filepath - the path to the file to be written
func GetWriterForFile(filepath string) Writer {

	for _, writer := range writers {
		if writer.CanWrite(filepath) {
			return writer
		}
	}

	return nil
}

//==========================================================================
// Private
//==========================================================================

var writers = []Writer{
	&JsonWriter{},
	&YamlWriter{},
	&TsWriter{},
	&PropertiesFileWriter{},
}
