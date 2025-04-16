package catalog

import (
	"log"
	"os"
	"os/exec"
	"testing"
)

func TestMain(m *testing.M) {
	// compile the testplugin binary used in some tests
	cmd := exec.Command("go", "build", "-buildvcs=false", "-o", "testpluginbinary", "./internal/testplugin")
	if output, err := cmd.CombinedOutput(); err != nil {
		log.Printf("output: %s", output)
		log.Fatalf("error: %v", err)
	}

	// run the tests
	code := m.Run()

	// remove the testplugin binary
	os.Remove("testpluginbinary")

	// exit with the code from the tests
	os.Exit(code)
}
