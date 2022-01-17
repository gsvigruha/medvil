package animation

type PersonMotion struct {
	LeftShoulder  [3]float64
	LeftHand      [8][3]float64
	LeftElbow     [8][3]float64
	RightShoulder [3]float64
	RightHand     [8][3]float64
	RightElbow    [8][3]float64
	LeftHip       [3]float64
	LeftKnee      [8][3]float64
	LeftFoot      [8][3]float64
	RightHip      [3]float64
	RightKnee     [8][3]float64
	RightFoot     [8][3]float64
}

var PersonMotionWalk = PersonMotion{
	LeftShoulder:  [3]float64{0, -25, -4},
	RightShoulder: [3]float64{0, -25, 4},
	LeftHip:       [3]float64{0, -15, -3},
	RightHip:      [3]float64{0, -15, 3},
	LeftElbow:     [8][3]float64{{0, -19, -4}, {1, -19.5, -4}, {2, -20, -4}, {1, -19.5, -4}, {0, -19, -4}, {-1, -19.5, -4}, {-2, -20, -4}, {-1, -19.5, -4}},
	RightElbow:    [8][3]float64{{0, -19, 4}, {-1, -19.5, 4}, {-2, -20, 4}, {-1, -19.5, 4}, {0, -19, 4}, {1, -19.5, 4}, {2, -20, 4}, {1, -19.5, 4}},
	LeftHand:      [8][3]float64{{0, -14, -4}, {2, -15, -4}, {4, -12, -4}, {2, -15, -4}, {0, -14, -4}, {-2, -15, -4}, {-4, -12, -4}, {-2, -15, -4}},
	RightHand:     [8][3]float64{{0, -14, 4}, {-2, -15, 4}, {-4, -12, 4}, {-2, -15, 4}, {0, -14, 4}, {2, -15, 4}, {4, -12, 4}, {2, -15, 4}},
	LeftKnee:      [8][3]float64{{0, -8, -3}, {1, -9, -3}, {2, -10, -3}, {1, -9, -3}, {0, -8, -3}, {-0.5, -8.5, -3}, {-1, -9, -3}, {-0.5, -8.5, -3}},
	RightKnee:     [8][3]float64{{0, -8, 3}, {-0.5, -8.5, 3}, {-1, -9, 3}, {-0.5, -8.5, 3}, {0, -8, 3}, {1, -9, 3}, {2, -10, 3}, {1, -9, 3}},
	LeftFoot:      [8][3]float64{{0, 0, -3}, {1, -1, -3}, {2, -2, -3}, {1, -1, -3}, {0, 0, -3}, {-1, 0, -3}, {-2, 1, -3}, {-1, 0, -3}},
	RightFoot:     [8][3]float64{{0, 0, 3}, {-1, 0, 3}, {-2, 1, 3}, {-1, 0, 3}, {0, 0, 3}, {1, -1, 3}, {2, -2, 3}, {1, -1, 3}},
}
