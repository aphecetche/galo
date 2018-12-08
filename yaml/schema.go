package yaml

type yaDigit struct {
	Deid   int
	Dsid   int
	Dsch   int
	Adc    int `yaml:"adc,omitempty"`
	Charge float32
}

type yaPixel struct {
	X  float32
	Y  float32
	DX float32
	DY float32
}

type yaStep struct {
	Pixels []yaPixel
	Ncalls int `yaml:"ncalls,omitempty"`
}

type yaDigitGroup struct {
	RefTime int // reference timestamp for the group digits
	Digits  []yaDigit
}

type yaPre struct {
	DigitGroup yaDigitGroup
}

type yaPos struct {
	X float32
	Y float32
}

type yaCluster struct {
	Pre    yaPre
	Pos    yaPos
	Charge float32
	Steps  []yaStep `yaml:"steps,omitempty"`
}

type yaDEClusters struct {
	DeID     int
	Clusters []yaCluster
}
