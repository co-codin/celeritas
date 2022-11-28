package celeritas

const version = "1.0.0"

type Celeritas struct {
	AppName string
	Debug   bool
	Version string
}



func (c *Celeritas) New(rootPath string) error {

}

func (c *Celeritas) Init(p initPaths) error {
    return nil
}