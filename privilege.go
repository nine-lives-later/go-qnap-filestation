package filestation

import "strconv"

// Privilege represents a linux file-system-level access privilege.
// It can be directly created as `Privilege(0644)` or from a string
// by using NewPrivilegeFromOctal().
type Privilege uint16

const (
	Privilege_Invalid Privilege = 0
)

type PrivilegeBits struct {
	OwnerRead    bool
	OwnerWrite   bool
	OwnerExecute bool
	GroupRead    bool
	GroupWrite   bool
	GroupExecute bool
	OtherRead    bool
	OtherWrite   bool
	OtherExecute bool
}

// NewPrivilegeFromOctal creates a new privilege by parsing
// a string like "644" or "0644".
// An empty or invalid string will return `Privilege_Invalid`.
func NewPrivilegeFromOctal(octalStr string) Privilege {
	if octalStr == "" {
		return Privilege_Invalid
	}
	r, err := strconv.ParseUint(octalStr, 8, 16)
	if err != nil {
		return Privilege_Invalid
	}
	return Privilege(r)
}

// String returns the octal representation of the privilege
// (without a leading zero), e.g. "644".
func (p Privilege) String() string {
	return strconv.FormatUint(uint64(p), 8)
}

// Bits returns the bitwise representation of the privilege.
func (p Privilege) Bits() PrivilegeBits {
	return PrivilegeBits{
		OwnerRead:    p&0400 != 0,
		OwnerWrite:   p&0200 != 0,
		OwnerExecute: p&0100 != 0,
		GroupRead:    p&0040 != 0,
		GroupWrite:   p&0020 != 0,
		GroupExecute: p&0010 != 0,
		OtherRead:    p&0004 != 0,
		OtherWrite:   p&0002 != 0,
		OtherExecute: p&0001 != 0,
	}
}
