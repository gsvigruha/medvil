package artifacts

type WorkProcess struct {
	Name    string
	Inputs  []*Artifacts
	Outputs []*Artifacts
}
