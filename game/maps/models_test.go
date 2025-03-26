package maps

import "testing"

func TestMap_Marshal(t *testing.T) {
	bytes, _ := Ctor(2, 2).Marshal()
	if string(bytes) != `{"width":2,"height":2,"values":[[0,0],[0,0]]}` {
		t.Fatal(string(bytes))
	}
}

func TestMap_Unmarshal(t *testing.T) {
	var m Map
	err := m.Unmarshal([]byte(`{"width":2,"height":2,"values":[[1,2],[3,4]]}`))
	if err != nil {
		t.Fatal(err)
	}
	if m.Width != 2 {
		t.Fatal(m.Width)
	}
}

func TestMap_At(t *testing.T) {
	m := Ctor(2, 2)
	*m.At(0, 0) = 1
	*m.At(0, 1) = 2
	if *m.At(0, 0) != 1 {
		t.Fatal(*m.At(0, 0))
	}
}
