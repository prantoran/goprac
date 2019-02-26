package ps2_test

import (
	"log"
	"testing"

	"github.com/prantoran/goprac/s2cells/ps2"
)

func TestS2CellID(t *testing.T) {
	// go test -v s2/s2_test.go -run TestS2CellID
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	orgLat, orgLon, dstLat, dstLon := 23.793751, 90.41130, 23.7851061163096, 90.40309242904186

	log.Printf("orgLat: %f, orgLon: %f, dstLat: %f, dstLon: %f\n",
		orgLat, orgLon, dstLat, dstLon)

	ps2.CellID(orgLat, orgLon)

}

func TestCellIDsLevelRange(t *testing.T) {
	// go test -v ps2/s2_test.go -run TestCellIDsLevelRange
	cellIDs, _ := ps2.CellIDsLevelRange(24.059, 90.55500000000006, 8, 13)

	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	// for _, u := range cellIDs {
	// 	log.Println(u)
	// }

	log.Println(len(cellIDs))

	// for
}
