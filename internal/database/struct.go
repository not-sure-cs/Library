package database 

type Linked struct {
	Author interface{} `json:"author"`
	Book interface{}  `json:"book"`
	Link   interface{} `json:"link"`	
}