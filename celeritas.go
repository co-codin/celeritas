package celeritas

const version = "1.0.0"

type Celeritas struct {
	AppName string
	Debug   bool
	Version string
}



func (c *Celeritas) New(rootPath string) error {
	return nil
}

func (c *Celeritas) Init(p initPaths) error {
    root := p.rootPath
	for _, path := range p.folderNames {
		// create folder if it doesnt exist
		err := c.CreateDirIfNotExists(root + "/" + path)
	    if err!= nil {
            return err
        }
    }
	return nil
}