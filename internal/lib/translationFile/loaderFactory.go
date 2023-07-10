package translationFile

//==========================================================================
// Public
//==========================================================================

// GetLoaderForFile returns the appropriate translation loader for the specified file
// #### params
// filepath - the path to the file to inspect
func GetLoaderForFile(filepath string) Loader {

	for _, loader := range loaders {
		if loader.CanLoad(filepath) {
			return loader
		}
	}

	return nil
}

//==========================================================================
// Private
//==========================================================================

var loaders = []Loader{
	&JsonLoader{},
	&YamlLoader{},
	&TsLoader{},
	&PropertiesFileLoader{},
}
