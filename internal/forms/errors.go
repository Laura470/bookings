package forms

type errors map[string][]string

//Add adds an error message for a given field, the field sono i name della form
func (e errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

// GEt returns the firs error message
func (e errors) Get(field string) string {
	es := e[field]
	if len(es) == 0 {
		return ""
	}
	//ritorna il primo campo
	return es[0]
}
