package photo

type EmbedResponse struct {
	Media Media `json:"shortcode_media"`
}

func (r EmbedResponse) getURLs() []string {
	if r.Media.IsVideo {
		return nil
	}
	s := r.Media.EdgeSidecar.getLowestImageURLs()
	if len(s) == 0 && r.Media.DisplayURL != "" {
		s = []string{r.Media.DisplayURL}
	}
	return s
}

// IsEmpty return if necessary script is absent
func (r EmbedResponse) IsEmpty() bool {
	return r.Media.ID == ""
}

type Media struct {
	ID          string      `json:"id"`
	TypeName    string      `json:"__typename"`
	IsVideo     bool        `json:"is_video"`
	DisplayURL  string      `json:"display_url"`
	EdgeSidecar EdgeSidecar `json:"edge_sidecar_to_children"`
}

type EdgeSidecar struct {
	Edges []Edge `json:"edges"`
}

func (es EdgeSidecar) getLowestImageURLs() (l []string) {
	for _, e := range es.Edges {
		if s, ok := e.Node.getLowestImageURL(); ok {
			l = append(l, s)
		}
	}
	return
}

type Edge struct {
	Node Node `json:"node"`
}

type Node struct {
	TypeName         string            `json:"__typename"`
	IsVideo          bool              `json:"is_video"`
	DisplayURL       string            `json:"display_url"`
	DisplayResources []DisplayResource `json:"display_resources"`
}

func (n Node) getLowestImageURL() (s string, ok bool) {
	if n.IsVideo {
		return "", false
	}
	s = n.DisplayURL
	index := -1
	minWidth := 10000
	for i, dr := range n.DisplayResources {
		if dr.Width < minWidth {
			minWidth = dr.Width
			index = i
		}
	}
	if index >= 0 {
		s = n.DisplayResources[index].Src
	}
	if s == "" {
		return s, false
	}
	return s, true
}

type DisplayResource struct {
	Width  int    `json:"config_width"`
	Height int    `json:"config_height"`
	Src    string `json:"src"`
}
