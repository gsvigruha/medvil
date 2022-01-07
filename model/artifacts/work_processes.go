package artifacts

type WorkProcess struct {
	Name string
	Inputs []*Artifact
	Outputs []*Artifact
}
