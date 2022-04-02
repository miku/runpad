package padutil

import (
	"context"
	"encoding/json"

	"github.com/FabianWe/etherpadlite-golang"
)

// Etherpad adds a few utility functions.
type Etherpad struct {
	etherpadlite.EtherpadLite
}

// PadTextAttribs is a response to a GetText response.
type PadTextAttribs struct {
	Code int64
	Data struct {
		Text struct {
			Attribs string `json:"attribs"`
			Text    string `json:"text"`
		} `json:"text"`
	}
	Message string
}

// PadText is a response to a GetText response.
type PadText struct {
	Code int64
	Data struct {
		Text string `json:"text"`
	}
	Message string
}

type PadSavedRevisions struct {
	Code int64
	Data struct {
		SavedRevisions []int64 `json:"savedRevisions"`
	}
	Message string
}

// Exists returns true, if a pad with a given id exists.
func (p *Etherpad) Exists(ctx context.Context, id string) (bool, error) {
	resp, err := p.GetText(ctx, id, 0)
	if err != nil {
		return false, err
	}
	return resp.Code == 0, nil
}

func (p *Etherpad) GetPadSavedRevisions(ctx context.Context, id string) (*PadSavedRevisions, error) {
	resp, err := p.ListSavedRevisions(ctx, id)
	if err != nil {
		return nil, err
	}
	b, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}
	var v PadSavedRevisions
	err = json.Unmarshal(b, &v)
	return &v, err
}

// GetPadText is a convenience function to get rid of type assertions.
func (p *Etherpad) GetPadText(ctx context.Context, id string) (*PadText, error) {
	resp, err := p.GetText(ctx, id, etherpadlite.OptionalParam)
	// {
	//   "Code": 0,
	//   "Message": "ok",
	//   "Data": {
	//     "text": {
	//       "attribs": "|1+s",
	//       "text": "Lorem ipsum dolor sit amet.\n"
	//     }
	//   }
	// }
	if err != nil {
		return nil, err
	}
	b, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}
	var v PadText
	err = json.Unmarshal(b, &v)
	return &v, err
}
