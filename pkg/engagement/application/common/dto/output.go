package dto

// Dummy ..
type Dummy struct {
	id string
}

//IsEntity ...
func (d Dummy) IsEntity() {}

// IsNode ..
func (d *Dummy) IsNode() {}

// SetID sets the trace's ID
func (d *Dummy) SetID(id string) {
	d.id = id
}
