package e2e

import (
	"testing"

	"github.com/samxsmith/snag"
)

func TestTwoFileTest(t *testing.T) {
	snag.Run("./files")

}
