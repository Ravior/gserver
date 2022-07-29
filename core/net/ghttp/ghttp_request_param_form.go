package ghttp

// GetForm retrieves and returns parameter <key> from form.
// It returns <def> if <key> does not exist in the form and <def> is given, or else it returns nil.
func (r *Request) GetForm(key string, def ...interface{}) interface{} {
	r.parseForm()
	if len(r.formMap) > 0 {
		if v, ok := r.formMap[key]; ok {
			return v
		}
	}
	if len(def) > 0 {
		return def[0]
	}
	return nil
}
