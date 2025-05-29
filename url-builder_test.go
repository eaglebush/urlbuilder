package urlbuilder

import "testing"

func TestBuild(t *testing.T) {
	ub := New(Host("localhost"), Port(5666))
	t.Logf("Simple: %s", ub.Build())

	ub1 := Clone(ub, Query("un", "zaldy.baguinon"))
	t.Logf("With Query String: %s", ub1.Build())

	ub2 := Clone(ub, Path("retrieve"), ID("key1"), Query("un", "zaldybaguinon"), Query("work", "ISD Manager"))
	t.Logf("With path and key: %s", ub2.Build())

	// Clone check
	t.Logf("Clone check: %s", ub.Build())

	ub3 := Clone(ub, Path("retrieve"), ID("key1"), Query("un", "zaldybaguinon"), Query("work", "ISD Manager"), UsrPwd("admin", "fantastic4"))
	t.Logf("With user, password, path and key: %s", ub3.Build())

	ub4 := NewSimpleUrl("localhost", "/path/")
	t.Logf("New Simple Url: %s", ub4.Build())

	ub5 := NewSimpleUrlWithID("localhost", "/path/", "12345")
	t.Logf("New Simple Url With ID: %s", ub5.Build())

	ub6 := ub5.Clone(Query("fn", "Elizalde"))
	t.Logf("New Simple Url With ID Cloned: %s", ub6.Build())

	// Clone check
	t.Logf("Clone check fr ub5: %s", ub5.Build())

	ub7 := ub5.Clone(Path("added-path"))
	t.Logf("Added path from ub5: %s", ub7.Build())
}
