package pdf

// Resolve resolves an indirect reference to its actual object.
// If the object is not a reference, it is returned as is.
func (r *Reader) Resolve(obj Object) (Object, error) {
	ref, ok := obj.(IndirectRef)
	if !ok {
		return obj, nil
	}

	return r.ReadObject(ref.ObjectNumber)
}
