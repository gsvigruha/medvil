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
	Tool          bool
}

var PersonMotionWalk = PersonMotion{
	LeftShoulder:  [3]float64{0, -26, -4},
	RightShoulder: [3]float64{0, -26, 4},
	LeftHip:       [3]float64{0, -16, -2},
	RightHip:      [3]float64{0, -16, 2},
	LeftElbow:     [8][3]float64{{0, -19, -4}, {1, -19.5, -4}, {2, -20, -4}, {1, -19.5, -4}, {0, -19, -4}, {-1, -19.5, -4}, {-2, -20, -4}, {-1, -19.5, -4}},
	RightElbow:    [8][3]float64{{0, -19, 4}, {-1, -19.5, 4}, {-2, -20, 4}, {-1, -19.5, 4}, {0, -19, 4}, {1, -19.5, 4}, {2, -20, 4}, {1, -19.5, 4}},
	LeftHand:      [8][3]float64{{0, -14, -4}, {2, -15, -4}, {4, -16, -4}, {2, -15, -4}, {0, -14, -4}, {-2, -15, -4}, {-4, -16, -4}, {-2, -15, -4}},
	RightHand:     [8][3]float64{{0, -14, 4}, {-2, -15, 4}, {-4, -16, 4}, {-2, -15, 4}, {0, -14, 4}, {2, -15, 4}, {4, -16, 4}, {2, -15, 4}},
	LeftKnee:      [8][3]float64{{0, -8, -2}, {1, -9, -2}, {2, -10, -2}, {1, -9, -2}, {0, -8, -2}, {-0.5, -8.5, -2}, {-1, -9, -2}, {-0.5, -8.5, -2}},
	RightKnee:     [8][3]float64{{0, -8, 2}, {-0.5, -8.5, 2}, {-1, -9, 2}, {-0.5, -8.5, 2}, {0, -8, 2}, {1, -9, 2}, {2, -10, 2}, {1, -9, 2}},
	LeftFoot:      [8][3]float64{{0, 0, -2}, {1, -1, -2}, {2, -2, -2}, {1, -1, -2}, {0, 0, -2}, {-1, 0, -2}, {-2, 1, -2}, {-1, 0, -2}},
	RightFoot:     [8][3]float64{{0, 0, 2}, {-1, 0, 2}, {-2, 1, 2}, {-1, 0, 2}, {0, 0, 2}, {1, -1, 2}, {2, -2, 2}, {1, -1, 2}},
	Tool:          false,
}

var PersonMotionFieldWork = PersonMotion{
	LeftShoulder:  [3]float64{0, -25, -4},
	RightShoulder: [3]float64{0, -25, 4},
	LeftHip:       [3]float64{0, -16, -2},
	RightHip:      [3]float64{0, -16, 2},
	LeftElbow:     [8][3]float64{{3, -20, -4}, {2.5, -19.5, -4}, {2, -19, -4}, {1.5, -18.5, -4}, {1, -18, -4}, {1.5, -18.5, -4}, {2, -19, -4}, {2.5, -19.5, -4}},
	RightElbow:    [8][3]float64{{-2, -19.5, 4}, {-2, -19.5, 4}, {-2, -19.5, 4}, {-2, -19.5, 4}, {-2, -19.5, 4}, {-2, -19.5, 4}, {-2, -19.5, 4}, {-2, -19.5, 4}},
	LeftHand:      [8][3]float64{{7, -17, -4}, {6.5, -16.5, -4}, {6, -16, -4}, {5.5, -15.5, -4}, {5, -15, -4}, {5.5, -15.5, -4}, {6, -16, -4}, {6.5, -16.5, -4}},
	RightHand:     [8][3]float64{{2, -19, 4}, {2, -18.5, 4}, {2, -18, 4}, {2, -17.5, 4}, {2, -17, 4}, {2, -17.5, 4}, {2, -18, 4}, {2, -18.5, 4}},
	LeftKnee:      [8][3]float64{{0, -8, -2}, {0, -8, -2}, {0, -8, -2}, {0, -8, -2}, {0, -8, -2}, {0, -8, -2}, {0, -8, -2}, {0, -8, -2}},
	RightKnee:     [8][3]float64{{0, -8, 2}, {0, -8, 2}, {0, -8, 2}, {0, -8, 2}, {0, -8, 2}, {0, -8, 2}, {0, -8, 2}, {0, -8, 2}},
	LeftFoot:      [8][3]float64{{0, 0, -2}, {0, 0, -2}, {0, 0, -2}, {0, 0, -2}, {0, 0, -2}, {0, 0, -2}, {0, 0, -2}, {0, 0, -2}},
	RightFoot:     [8][3]float64{{0, 0, 2}, {0, 0, 2}, {0, 0, 2}, {0, 0, 2}, {0, 0, 2}, {0, 0, 2}, {0, 0, 2}, {0, 0, 2}},
	Tool:          true,
}

var PersonMotionBuild = PersonMotion{
	LeftShoulder:  [3]float64{0, -25, -4},
	RightShoulder: [3]float64{0, -25, 4},
	LeftHip:       [3]float64{0, -16, -2},
	RightHip:      [3]float64{0, -16, 2},
	LeftElbow:     [8][3]float64{{0, -20, -6}, {0, -20, -6}, {0, -20, -6}, {0, -20, -6}, {0, -20, -6}, {0, -20, -6}, {0, -20, -6}, {0, -20, -6}},
	RightElbow:    [8][3]float64{{1, -21, 6}, {1, -21, 6}, {1, -21, 6}, {1, -21, 6}, {1, -21, 6}, {1, -21, 6}, {1, -21, 6}, {1, -21, 6}},
	LeftHand:      [8][3]float64{{0, -25, -8}, {0, -25, -8}, {0, -25, -8}, {0, -25, -8}, {0, -25, -8}, {0, -25, -8}, {0, -25, -8}, {0, -25, -8}},
	RightHand:     [8][3]float64{{1, -26, 6}, {0, -26, 6}, {-1, -27, 6}, {-2, -27, 6}, {-3, -28, 6}, {-2, -27, 6}, {-1, -27, 6}, {0, -26, 6}},
	LeftKnee:      [8][3]float64{{0, -8, -2}, {0, -8, -2}, {0, -8, -2}, {0, -8, -2}, {0, -8, -2}, {0, -8, -2}, {0, -8, -2}, {0, -8, -2}},
	RightKnee:     [8][3]float64{{0, -8, 2}, {0, -8, 2}, {0, -8, 2}, {0, -8, 2}, {0, -8, 2}, {0, -8, 2}, {0, -8, 2}, {0, -8, 2}},
	LeftFoot:      [8][3]float64{{0, 0, -2}, {0, 0, -2}, {0, 0, -2}, {0, 0, -2}, {0, 0, -2}, {0, 0, -2}, {0, 0, -2}, {0, 0, -2}},
	RightFoot:     [8][3]float64{{0, 0, 2}, {0, 0, 2}, {0, 0, 2}, {0, 0, 2}, {0, 0, 2}, {0, 0, 2}, {0, 0, 2}, {0, 0, 2}},
	Tool:          false,
}

var PersonMotionMine = PersonMotion{
	LeftShoulder:  [3]float64{0, -25, -4},
	RightShoulder: [3]float64{0, -25, 4},
	LeftHip:       [3]float64{0, -16, -2},
	RightHip:      [3]float64{0, -16, 2},
	LeftElbow:     [8][3]float64{{3, -17, -4}, {2.5, -16.5, -4}, {2, -16, -4}, {1.5, -15.5, -4}, {1, -15, -4}, {1.5, -15.5, -4}, {2, -16, -4}, {2.5, -16.5, -4}},
	RightElbow:    [8][3]float64{{-2, -19.5, 4}, {-2, -19.5, 4}, {-2, -19.5, 4}, {-2, -19.5, 4}, {-2, -19.5, 4}, {-2, -19.5, 4}, {-2, -19.5, 4}, {-2, -19.5, 4}},
	LeftHand:      [8][3]float64{{7, -19, -4}, {6.5, -18, -4}, {6, -17, -4}, {5.5, -16, -4}, {5, -15, -4}, {5.5, -16, -4}, {6, -17, -4}, {6.5, -18, -4}},
	RightHand:     [8][3]float64{{2, -19, 4}, {2, -18.5, 4}, {2, -18, 4}, {2, -17.5, 4}, {2, -17, 4}, {2, -17.5, 4}, {2, -18, 4}, {2, -18.5, 4}},
	LeftKnee:      [8][3]float64{{0, -8, -2}, {0, -8, -2}, {0, -8, -2}, {0, -8, -2}, {0, -8, -2}, {0, -8, -2}, {0, -8, -2}, {0, -8, -2}},
	RightKnee:     [8][3]float64{{0, -8, 2}, {0, -8, 2}, {0, -8, 2}, {0, -8, 2}, {0, -8, 2}, {0, -8, 2}, {0, -8, 2}, {0, -8, 2}},
	LeftFoot:      [8][3]float64{{0, 0, -2}, {0, 0, -2}, {0, 0, -2}, {0, 0, -2}, {0, 0, -2}, {0, 0, -2}, {0, 0, -2}, {0, 0, -2}},
	RightFoot:     [8][3]float64{{0, 0, 2}, {0, 0, 2}, {0, 0, 2}, {0, 0, 2}, {0, 0, 2}, {0, 0, 2}, {0, 0, 2}, {0, 0, 2}},
	Tool:          true,
}

var PersonMotionCut = PersonMotion{
	LeftShoulder:  [3]float64{0, -25, -4},
	RightShoulder: [3]float64{0, -25, 4},
	LeftHip:       [3]float64{0, -16, -2},
	RightHip:      [3]float64{0, -16, 2},
	LeftElbow:     [8][3]float64{{1, -20, -2}, {1, -20, -1.5}, {1, -20, -1}, {1, -20, -0.5}, {1, -20, 0}, {1, -20, -0.5}, {1, -20, -1}, {1, -20, -1.5}},
	RightElbow:    [8][3]float64{{-2, -20, 4}, {-2, -20, 4}, {-2, -20, 4}, {-2, -20, 4}, {-2, -20, 4}, {-2, -20, 4}, {-2, -20, 4}, {-2, -20, 4}},
	LeftHand:      [8][3]float64{{7, -20, 2}, {6.5, -20, 3}, {6, -20, 4}, {5.5, -20, 5}, {5, -20, 6}, {5.5, -20, 5}, {6, -20, 4}, {6.5, -20, 3}},
	RightHand:     [8][3]float64{{3, -20, 4}, {3, -20, 4.5}, {2.5, -20, 5}, {2.5, -20, 5.5}, {2, -20, 6}, {2.5, -20, 5.5}, {2.5, -20, 5}, {3, -20, 4.5}},
	LeftKnee:      [8][3]float64{{0, -8, -2}, {0, -8, -2}, {0, -8, -2}, {0, -8, -2}, {0, -8, -2}, {0, -8, -2}, {0, -8, -2}, {0, -8, -2}},
	RightKnee:     [8][3]float64{{0, -8, 2}, {0, -8, 2}, {0, -8, 2}, {0, -8, 2}, {0, -8, 2}, {0, -8, 2}, {0, -8, 2}, {0, -8, 2}},
	LeftFoot:      [8][3]float64{{0, 0, -2}, {0, 0, -2}, {0, 0, -2}, {0, 0, -2}, {0, 0, -2}, {0, 0, -2}, {0, 0, -2}, {0, 0, -2}},
	RightFoot:     [8][3]float64{{0, 0, 2}, {0, 0, 2}, {0, 0, 2}, {0, 0, 2}, {0, 0, 2}, {0, 0, 2}, {0, 0, 2}, {0, 0, 2}},
	Tool:          true,
}

var PersonMotionPaddle = PersonMotion{
	LeftShoulder:  [3]float64{0, -25, -4},
	RightShoulder: [3]float64{0, -25, 4},
	LeftHip:       [3]float64{0, -16, -2},
	RightHip:      [3]float64{0, -16, 2},
	LeftElbow:     [8][3]float64{{0, -20, -6}, {0.5, -20, -6}, {1, -20, -6}, {1.5, -20, -6}, {2, -20, -6}, {1.5, -20, -6}, {1, -20, -6}, {0.5, -20, -6}},
	RightElbow:    [8][3]float64{{0, -20, 6}, {0.5, -20, 6}, {1, -20, 6}, {1.5, -20, 6}, {2, -20, 6}, {1.5, -20, 6}, {1, -20, 6}, {0.5, -20, 6}},
	LeftHand:      [8][3]float64{{0, -17, -7}, {1, -17, -7}, {2, -17, -7}, {3, -17, -7}, {4, -17, -7}, {3, -17, -7}, {2, -17, -7}, {1, -17, -7}},
	RightHand:     [8][3]float64{{0, -17, 7}, {1, -17, 7}, {2, -17, 7}, {3, -17, 7}, {4, -17, 7}, {3, -17, 7}, {2, -17, 7}, {1, -17, 7}},
	// Legs are not rendered for paddling
	LeftKnee:  [8][3]float64{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}, {0, 0, 0}, {0, 0, 0}, {0, 0, 0}, {0, 0, 0}, {0, 0, 0}},
	RightKnee: [8][3]float64{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}, {0, 0, 0}, {0, 0, 0}, {0, 0, 0}, {0, 0, 0}, {0, 0, 0}},
	LeftFoot:  [8][3]float64{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}, {0, 0, 0}, {0, 0, 0}, {0, 0, 0}, {0, 0, 0}, {0, 0, 0}},
	RightFoot: [8][3]float64{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}, {0, 0, 0}, {0, 0, 0}, {0, 0, 0}, {0, 0, 0}, {0, 0, 0}},
	Tool:      false,
}
