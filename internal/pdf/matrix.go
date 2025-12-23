package pdf

// Matrix represents a 3x3 affine transformation matrix [a b c d e f]
// a b 0
// c d 0
// e f 1
type Matrix [6]float64

// IdentityMatrix returns [1 0 0 1 0 0]
func IdentityMatrix() Matrix {
	return Matrix{1, 0, 0, 1, 0, 0}
}

// Multiply multiplies two matrices: m x other
func (m Matrix) Multiply(other Matrix) Matrix {
	// [a1 b1 0]   [a2 b2 0]
	// [c1 d1 0] x [c2 d2 0]
	// [e1 f1 1]   [e2 f2 1]

	return Matrix{
		m[0]*other[0] + m[1]*other[2],            // a
		m[0]*other[1] + m[1]*other[3],            // b
		m[2]*other[0] + m[3]*other[2],            // c
		m[2]*other[1] + m[3]*other[3],            // d
		m[4]*other[0] + m[5]*other[2] + other[4], // e (tx)
		m[4]*other[1] + m[5]*other[3] + other[5], // f (ty)
	}
}

// Transform transforms a point (x, y)
func (m Matrix) Transform(x, y float64) (float64, float64) {
	// [x y 1] * m
	tx := x*m[0] + y*m[2] + m[4]
	ty := x*m[1] + y*m[3] + m[5]
	return tx, ty
}
