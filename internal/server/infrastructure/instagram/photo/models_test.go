package photo

import (
	"reflect"
	"testing"
)

func TestNode_getLowestImageURL(t *testing.T) {
	type fields struct {
		TypeName         string
		IsVideo          bool
		DisplayURL       string
		DisplayResources []DisplayResource
	}
	tests := []struct {
		name   string
		fields fields
		wantS  string
		wantOk bool
	}{
		{
			name: "valid_display",
			fields: fields{
				IsVideo:          false,
				DisplayURL:       "abc",
				DisplayResources: nil,
			},
			wantS: "abc",
			wantOk: true,
		},
		{
			name: "valid",
			fields: fields{
				IsVideo:          false,
				DisplayURL:       "abc",
				DisplayResources: []DisplayResource{
					{
						Width:  5,
						Height: 10,
						Src:    "qwe",
					},
					{
						Width:  7,
						Height: 10,
						Src:    "zxc",
					},
				},
			},
			wantS: "qwe",
			wantOk: true,
		},
		{
			name: "unvalid video",
			fields: fields{
				IsVideo:          true,
				DisplayURL:       "abc",
				DisplayResources: []DisplayResource{
					{
						Width:  5,
						Height: 10,
						Src:    "qwe",
					},
					{
						Width:  7,
						Height: 10,
						Src:    "zxc",
					},
				},
			},
			wantS: "",
			wantOk: false,
		},
		{
			name: "unvalid",
			fields: fields{
				IsVideo:          false,
				DisplayURL:       "",
				DisplayResources: []DisplayResource{},
			},
			wantS: "",
			wantOk: false,
		},
		{
			name: "valid_2",
			fields: fields{
				IsVideo:          false,
				DisplayURL:       "a",
				DisplayResources: []DisplayResource{},
			},
			wantS: "a",
			wantOk: true,
		},
	}
		for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := Node{
				TypeName:         tt.fields.TypeName,
				IsVideo:          tt.fields.IsVideo,
				DisplayURL:       tt.fields.DisplayURL,
				DisplayResources: tt.fields.DisplayResources,
			}
			gotS, gotOk := n.getLowestImageURL()
			if gotS != tt.wantS {
				t.Errorf("getLowestImageURL() gotS = %v, want %v", gotS, tt.wantS)
			}
			if gotOk != tt.wantOk {
				t.Errorf("getLowestImageURL() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func TestEdgeSidecar_getLowestImageURLs(t *testing.T) {
	type fields struct {
		Edges []Edge
	}
	tests := []struct {
		name   string
		fields fields
		wantL  []string
	}{
		{
			name: "valid",
			fields: fields{
				Edges: []Edge{
					{
						Node: Node{
							IsVideo:          false,
							DisplayURL:       "abc",
							DisplayResources: nil,
						},
					},
					{
						Node: Node{
							IsVideo:          false,
							DisplayURL:       "abc",
							DisplayResources: []DisplayResource{
								{
									Width:  5,
									Height: 10,
									Src:    "qwe",
								},
								{
									Width:  7,
									Height: 10,
									Src:    "zxc",
								},
							},
						},
					},
					{
						Node: Node{
							IsVideo:          true,
							DisplayURL:       "abc",
							DisplayResources: []DisplayResource{
								{
									Width:  5,
									Height: 10,
									Src:    "qwe",
								},
								{
									Width:  7,
									Height: 10,
									Src:    "zxc",
								},
							},
						},
					},
					{
						Node: Node{
							IsVideo:          false,
							DisplayURL:       "",
							DisplayResources: []DisplayResource{},
						},
					},
				},
			},
			wantL: []string{"abc", "qwe"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			es := EdgeSidecar{
				Edges: tt.fields.Edges,
			}
			if gotL := es.getLowestImageURLs(); !reflect.DeepEqual(gotL, tt.wantL) {
				t.Errorf("getLowestImageURLs() = %v, want %v", gotL, tt.wantL)
			}
		})
	}
}

func TestEmbedResponse_getURLs(t *testing.T) {
	type fields struct {
		Media Media
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			name: "non valid",
			fields: fields{
				Media: Media{
					IsVideo:     true,
					EdgeSidecar: EdgeSidecar{},
				},
			},
			want: nil,
		},
		{
			name: "valid",
			fields: fields{
				Media: Media{
					IsVideo:     false,
					EdgeSidecar: EdgeSidecar{
						Edges: []Edge{
							{
								Node: Node{
									IsVideo:          false,
									DisplayURL:       "abc",
									DisplayResources: nil,
								},
							},
							{
								Node: Node{
									IsVideo:          false,
									DisplayURL:       "abc",
									DisplayResources: []DisplayResource{
										{
											Width:  5,
											Height: 10,
											Src:    "qwe",
										},
										{
											Width:  7,
											Height: 10,
											Src:    "zxc",
										},
									},
								},
							},
							{
								Node: Node{
									IsVideo:          true,
									DisplayURL:       "abc",
									DisplayResources: []DisplayResource{
										{
											Width:  5,
											Height: 10,
											Src:    "qwe",
										},
										{
											Width:  7,
											Height: 10,
											Src:    "zxc",
										},
									},
								},
							},
							{
								Node: Node{
									IsVideo:          false,
									DisplayURL:       "",
									DisplayResources: []DisplayResource{},
								},
							},
						},
					},
				},
			},
			want: []string{"abc", "qwe"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := EmbedResponse{
				Media: tt.fields.Media,
			}
			if got := r.getURLs(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getURLs() = %v, want %v", got, tt.want)
			}
		})
	}
}
