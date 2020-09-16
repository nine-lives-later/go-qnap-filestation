package filestation

import (
	"fmt"
	"testing"
)

func TestPrivilegeFromOctal(t *testing.T) {
	p := NewPrivilegeFromOctal("644")
	if p == Privilege_Invalid {
		t.Fatal("privilege is invalid")
	}
	if uint16(p) != 0644 {
		t.Fatalf("wrong privilege value: %v", p)
	}
	if p.String() != "644" {
		t.Fatalf("wrong privilege string: %v", p)
	}
	if fmt.Sprintf("%v", p) != "644" {
		t.Fatalf("wrong privilege formatted string: %v", p)
	}
}

func TestPrivilegeFromInteger(t *testing.T) {
	p := Privilege(0644)
	if uint16(p) != 0644 {
		t.Fatalf("wrong privilege value: %v", p)
	}
	if p.String() != "644" {
		t.Fatalf("wrong privilege string: %v", p)
	}
}

func TestPrivilegeFromOctalWithLeadingZero(t *testing.T) {
	p := NewPrivilegeFromOctal("0644")
	if p == Privilege_Invalid {
		t.Fatal("privilege is invalid")
	}
	if uint16(p) != 0644 {
		t.Fatalf("wrong privilege value: %v", p)
	}
}

func TestPrivilegeBits(t *testing.T) {
	p := Privilege(0764)
	b := p.Bits()
	if !b.OwnerExecute || !b.OwnerWrite || !b.OwnerRead ||
		b.GroupExecute || !b.GroupWrite || !b.GroupRead ||
		b.OtherExecute || b.OtherWrite || !b.OtherRead {
		t.Fatalf("wrong privilege bits: %+v", b)
	}
}
