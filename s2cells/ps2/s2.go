package ps2

import (
	"errors"
	"log"

	"github.com/golang-collections/collections/stack"
	"github.com/golang/geo/s2"
)

// CellID takes lat, lng in float64 and returns s2cell ID
func CellID(lat float64, lng float64) uint64 {

	latlng := s2.LatLngFromDegrees(lat, lng)

	latA := latlng.Lat

	log.Println(
		"lat degrees: ", latA.Degrees(),
		"\nradians: ", latA.Radians(),
		"\nE7: ", latA.E7(),
		"\nAbs: ", latA.Abs(),
		"\nNormalized: ", latA.Normalized(),
		"\nString: ", latA.String(),
	)

	point := s2.PointFromLatLng(latlng)

	log.Println("point x: ", point.X, " y: ", point.Y, " z: ", point.Z)

	cid := s2.CellIDFromLatLng(latlng)

	log.Println(
		"cellID: ", cid,
		"\nToToken: ", cid.ToToken(),
		"\nIsValid: ", cid.IsValid(),
		"\nFace: ", cid.Face(),
		"\nPos: ", cid.Pos(),
		"\nLevel: ", cid.Level(),
		"\nIsLeaf: ", cid.IsLeaf(),
		"\nParent_level13: ", cid.Parent(13),
		"\nparent _level13's calculated level: ", cid.Parent(13).Level(),
		"\nString: ", cid.String(),
		"\nPoint: ", cid.Point(),
		"\nLatLng: ", cid.LatLng(),
		"\nfromtoken:", s2.CellIDFromToken(cid.ToToken()),
	)

	parent := cid.Parent(29)
	log.Println("parent contains cid:", parent.Contains(cid))

	return uint64(cid)

}

func Region(lat float64, lng float64) s2.Rect {
	ll := s2.LatLngFromDegrees(lat, lng)
	rect := s2.RectFromLatLng(ll)

	return rect
}

var st *stack.Stack

func CellIDsLevelRange(lat, lng float64, levelLo, levelHi int) ([]s2.CellID, error) {

	if levelLo > levelHi || levelLo < 0 || levelHi > 30 {
		return nil, errors.New("s2 cell levels not valid")
	}

	ll := s2.LatLngFromDegrees(lat, lng)

	cid := s2.CellIDFromLatLng(ll)

	levelLoCellID := cid.Parent(levelLo)

	c := s2.CellFromCellID(levelLoCellID)

	log.Println(c)

	st = stack.New()
	st.Push(c)

	ret := []s2.CellID{}

	for st.Len() > 0 {
		curCell := st.Pop().(s2.Cell)

		if curCell.Level() == levelHi {
			ret = append(ret, curCell.ID())
			continue
		}

		children, _ := curCell.Children()

		for _, child := range children {
			st.Push(child)
		}
	}

	return ret, nil
}
