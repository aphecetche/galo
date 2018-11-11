package cmd

import "github.com/aphecetche/pigiron/mapping"

type SegPair struct {
	Bending, NonBending mapping.Segmentation
}

var segmentations map[int]SegPair

func segmentation(deid int, bending bool) mapping.Segmentation {
	seg := segmentations[deid]
	if seg.Bending == nil {
		segmentations[deid] = SegPair{
			Bending:    mapping.NewSegmentation(deid, true),
			NonBending: mapping.NewSegmentation(deid, false),
		}
		seg = segmentations[deid]
	}
	if bending {
		return seg.Bending
	}
	return seg.NonBending
}
