package translationFile

import (
	_ "embed"
	"testing"
)

func TestTsLoader_CanLoad(t *testing.T) {
	tsLoader := TsLoader{}

	if !tsLoader.CanLoad("foobar.ts") || !tsLoader.CanLoad("/var/lib/fiz.js") {
		t.Fatalf("TS Loader said it can't load files that it should be able to load")
	}

	if tsLoader.CanLoad("text.csv") || tsLoader.CanLoad("../../foobar.jpg") {
		t.Fatalf("TS Loader said it could load files that it cant")
	}

}
