package model

import (
	"testing"

	"ymir/pkg/db"
)

func TestModelService_ListModels(t *testing.T) {
	tests := []struct {
		name    string
		want    int
		wantErr bool
	}{
		{"Test1", 2, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := NewModelService().(ModelService)
			got, err := ms.ListModels()
			if err != nil {
				t.Error(err)
				return
			}
			if tt.want != len(got) {
				t.Errorf("ListModels() got = %v, want %v", len(got), tt.want)
			}

			//bytes, err := json.MarshalIndent(got, "", "\t")
			//fmt.Println(string(bytes))
		})
	}
}

func TestModelService_GetModelsByTag(t *testing.T) {
	type fields struct {
		name      string
		DataStore *db.DB
	}
	type args struct {
		tag string
	}
	tests := []struct {
		name    string
		tag     string
		want    int
		wantErr bool
	}{
		{"tag3", "tag3", 1, true},
		{"tag2", "tag2", 2, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := NewModelService().(ModelService)
			got, err := ms.GetModelsByTag(tt.tag)
			if err != nil {
				t.Error(err)
				return
			}
			if tt.want != len(got) {
				t.Errorf("ListModels() got = %v, want %v", len(got), tt.want)
			}

			//bytes, err := json.MarshalIndent(got, "", "\t")
			//fmt.Println(string(bytes))
		})
	}
}
