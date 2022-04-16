package navigation

const DirectionN = 0
const DirectionE = 1
const DirectionS = 2
const DirectionW = 3

const DirectionNone = 255

var DirectionXY = [4][2]int{{0, -1}, {1, 0}, {0, 1}, {-1, 0}}
var DirectionDiagonalXY = [4][2]int{{1, -1}, {1, 1}, {-1, 1}, {-1, -1}}
