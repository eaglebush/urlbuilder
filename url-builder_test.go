package urlbuilder

import "testing"

func TestBuild(t *testing.T) {
	ub := New(Port(5666))
	t.Logf("No host: %s", ub.Build())

	ub0 := New(Host("localhost.com"), Port(5666))
	t.Logf("Simple: %s", ub0.Build())

	ub1 := Clone(ub0, Query("un", "zaldy.baguinon"))
	t.Logf("With Query String: %s", ub1.Build())

	ub2 := Clone(ub0, Path("retrieve"), ID("key1"), Query("un", "zaldy baguinon"), Query("work", "ISD Manager"))
	t.Logf("With path, key and query with spaces: %s", ub2.Build())

	// Clone check
	t.Logf("Clone check: %s", ub0.Build())

	ub3 := Clone(ub0, Path("retrieve"), ID("key1"), Query("un", "zaldybaguinon"), Query("work", "ISD Manager"), UsrPwd("admin", "fantastic4"))
	t.Logf("With user, password, path and key: %s", ub3.Build())

	ub4 := NewUrlWithPath("localhost", "/path/")
	t.Logf("New Simple Url: %s", ub4.Build())

	ub5 := NewUrlWithID("localhost", "/path/", "12345", Path("udoms"))
	t.Logf("New Simple Url With ID: %s", ub5.Build())

	ub6 := ub5.Clone(Query("fn", "Elizalde"))
	t.Logf("New Simple Url With ID Cloned: %s", ub6.Build())

	// Clone check
	t.Logf("Clone check fr ub5: %s", ub5.Build())

	ub7 := ub5.Clone(Path("added-path"))
	t.Logf("Added path from ub5: %s", ub7.Build())

	ub8 := New(Host("localhost"), ID(12345))
	t.Logf("Host plus key: %s", ub8.Build())

	t.Logf("Inline clone build: %s", ub8.Clone(Query("yes", "no")).Build())

	ub9 := New(Host("https://www.facebook.com:1500/ui/from/u"), Path("ever"), Query("open", 1), ID(324))
	t.Logf("Literal host: %s", ub9.Build())

	ub10 := NewUrlWithPath("localhost:3000", "", Path("/grpperm/"))
	t.Logf("Host with port and blank first path: %s", ub10.Build())

}

func TestSingle(t *testing.T) {
	// ub8 := NewUrlWithPath("localhost:3000", "", Path("/grpperm/"))
	// t.Logf("Host with port and blank first path: %s", ub8.Build())

	// ub9 := New(Host("https://www.facebook.com:1500/ui/from/u"), Path("ever"), Query("open", 1), ID(324))
	// t.Logf("Literal host: %s", ub9.Build())

	// ub10 := NewUrlWithPath("http://localhost:3000", "", Path("/grpperm/"))
	// t.Logf("Host with scheme and port and blank first path: %s", ub10.Build())

	// ub11 := NewUrlWithPath("localhost:3000", "", Path("/grpperm/"), Sch("http"))
	// t.Logf("Host with scheme and port and blank first path: %s", ub11.Build())

	// ub12 := NewUrlWithPath("http://localhost:3000/", "", Path("/grpperm/"), Sch("http"))
	// t.Logf("Host with scheme and port and blank first path: %s", ub12.Build())

	ub13 := NewUrl("https://appcore-test.vdimdci.com.ph/api", Path("user"), Path("info"), Query("userid", "zaldy.baguinon"))
	t.Logf("Host with scheme and port and blank first path: %s", ub13.Build())
}

func TestQueryStringBuild(t *testing.T) {
	qs := NewQueryString(
		QModeArray,
		Nv("last", "yes"),
		Nv("first", "no"),
		Nv("country", "Philippines"),
	)
	t.Logf("Query string build: %s", qs.Build())
}

func BenchmarkSimpleURL(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewUrlWithPath("example.com", "api/v1/users").Build()
	}
}

func BenchmarkURLWithID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewUrlWithID("example.com", "api/v1/users", 12345).Build()
	}
}

func BenchmarkURLWithQueryParams(b *testing.B) {
	for i := 0; i < b.N; i++ {
		New(
			Host("example.com"),
			Path("api/v1/search"),
			Query("q", "golang"),
			Query("page", 2),
			Query("sort", "desc"),
			Mode(QModeLast),
		).Build()
	}
}

func BenchmarkComplexURL(b *testing.B) {
	for i := 0; i < b.N; i++ {
		New(
			Sch("http"),
			Host("example.com"),
			UsrPwd("user", "pass"),
			Path("api"),
			Path("v1"),
			Path("resource"),
			ID("abcd-1234"),
			Query("filter", "active"),
			Query("limit", 50),
			Query("offset", 100),
			Mode(QModeArray),
			Port(8080),
			Frag("section"),
		).Build()
	}
}

// func BenchmarkComplexURL2(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		New(
// 			Sch("http"),
// 			Host("example.com"),
// 			UsrPwd("user", "pass"),
// 			Path("api"),
// 			Path("v1"),
// 			Path("resource"),
// 			ID("abcd-1234"),
// 			Query("filter", "active"),
// 			Query("limit", 50),
// 			Query("offset", 100),
// 			Mode(QModeArray),
// 			Port(8080),
// 			Fragment("section"),
// 		).Build2()
// 	}
// }
