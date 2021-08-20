package dto

import (
	"testing"
)

func TestDummy_IsEntity(t *testing.T) {
	type fields struct {
		id string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "happy case",
			fields: fields{
				id: "test id",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := Dummy{
				id: tt.fields.id,
			}
			d.IsEntity()
		})
	}
}

func TestDummy_IsNode(t *testing.T) {
	type fields struct {
		id string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "default case",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Dummy{
				id: tt.fields.id,
			}
			d.IsNode()
		})
	}
}

func TestDummy_SetID(t *testing.T) {
	type fields struct {
		id string
	}
	type args struct {
		id string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "good case",
			args: args{
				id: "an ID",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Dummy{
				id: tt.fields.id,
			}
			d.SetID(tt.args.id)
		})
	}
}
