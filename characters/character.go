package characters

type Character struct {
	Name   string
	Family string
}

type Characters map[string]Character
type Data struct {
	characters map[string]Character
}
