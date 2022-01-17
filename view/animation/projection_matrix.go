package animation

type ProjectionMatrix struct {
	XX float64
	XY float64
	XZ float64
	YX float64
	YY float64
	YZ float64
}

var ProjectionMatrixNE = ProjectionMatrix{
	XX: -0.2,
	XY: 0.0,
	XZ: -0.8,
	YX: 0.1,
	YY: 1.0,
	YZ: 0.1,
}

var ProjectionMatrixSE = ProjectionMatrix{
	XX: -0.2,
	XY: 0.0,
	XZ: -0.8,
	YX: 0.1,
	YY: 1.0,
	YZ: 0.1,
}

var ProjectionMatrixSW = ProjectionMatrix{
	XX: 0.2,
	XY: 0.0,
	XZ: 0.8,
	YX: -0.1,
	YY: 1.0,
	YZ: -0.1,
}

var ProjectionMatrixNW = ProjectionMatrix{
	XX: 0.2,
	XY: 0.0,
	XZ: 0.8,
	YX: -0.1,
	YY: 1.0,
	YZ: -0.1,
}
