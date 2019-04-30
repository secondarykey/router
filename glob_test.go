package router_test

import (
	"log"
	"testing"

	"github.com/secondarykey/router"
)

func TestPart(t *testing.T) {

	g, err := router.Compile("t[e]st")
	if err != nil {
		t.Errorf(`"t[e]st" error is not nil`)
	}
	if len(g) != 3 {
		t.Errorf(`"t[e]st" length 3"`)
	}
	log.Println(g)

	g, err = router.Compile("t[a-z]+st")
	if err != nil {
		t.Errorf(`"t[a-z]+st" error is not nil`)
	}
	if len(g) != 3 {
		t.Errorf(`"t[a-z]+st" length 3"`)
	}
	log.Println(g)

	g, err = router.Compile("t*st")
	if len(g) != 3 {
		t.Errorf(`"t*st" length 3 but %d"`, len(g))
	}
	log.Println(g)

	g, err = router.Compile("*t*st*")
	if len(g) != 5 {
		t.Errorf(`"t*st" length 5 but %d"`, len(g))
	}
	log.Println(g)

	g, err = router.Compile("*t?st*")
	if len(g) != 5 {
		t.Errorf(`"*t?st*" length 5 but %d"`, len(g))
	}
	log.Println(g)

	g, err = router.Compile("*t?[a-z]+st*")
	if len(g) != 6 {
		t.Errorf(`"*t?[a-z]st*" length 6 but %d"`, len(g))
	}
	log.Println(g)
}

func TestGlob(t *testing.T) {
	/*
		g, err := router.Compile("t[e]st")
		if err != nil {
			t.Errorf(`"t[e]st" is not nil`)
		}
	*/

}
