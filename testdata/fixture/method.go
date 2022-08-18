package fixture

import "fmt"

type Ob struct {
	name string
}

func (ob Ob) Name() string {
	return ob.name
}

func (ob *Ob) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`{"name": %q}`, ob.Name)), nil
}
