package data_structure

type Identifier interface {
	Id() string
}

type IdentifierSet struct {
	Ids map[string]Identifier
}

func (s IdentifierSet) Add(id Identifier) {
	s.Ids[id.Id()] = id
}

func (s *IdentifierSet) Remove(id string) {
	delete(s.Ids, id)
}

func (s *IdentifierSet) Contains(id string) bool {
	_, ok := s.Ids[id]
	return ok
}

func (s *IdentifierSet) Size() int {
	return len(s.Ids)
}
