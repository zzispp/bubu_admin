package v1

type CreateRoleRequest struct {
	Code        string   `protobuf:"bytes,1,opt,name=code,proto3" json:"code,omitempty"`               // Code of role (unique)
	Name        string   `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`               // Display name of role
	Description string   `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"` // Details about role
	Sequence    int32    `protobuf:"varint,4,opt,name=sequence,proto3" json:"sequence,omitempty"`      // Sequence for sorting
	Status      string   `protobuf:"bytes,5,opt,name=status,proto3" json:"status,omitempty"`           // Status of role (enabled, disabled)
	Menus       []string `protobuf:"bytes,6,rep,name=menus,proto3" json:"menus,omitempty"`             // Role menu list
}