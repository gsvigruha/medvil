package animation

type ProjectionMatrix struct {
	XX float64
	XY float64
	XZ float64
	YX float64
	YY float64
	YZ float64
}

var ProjectionMatrices = [4]ProjectionMatrix{
	ProjectionMatrix{
		XX: 0.83,
		XY: 0.0,
		XZ: 0.83,
		YX: 0.55,
		YY: 1.0,
		YZ: -0.55,
	},
	ProjectionMatrix{
		XX: 0.83,
		XY: 0.0,
		XZ: -0.83,
		YX: -0.55,
		YY: 1.0,
		YZ: -0.55,
	},
	ProjectionMatrix{
		XX: -0.83,
		XY: 0.0,
		XZ: -0.83,
		YX: -0.55,
		YY: 1.0,
		YZ: 0.55,
	},
	ProjectionMatrix{
		XX: -0.83,
		XY: 0.0,
		XZ: 0.83,
		YX: 0.55,
		YY: 1.0,
		YZ: 0.55,
	},
}
