package model

// core.Configではなく今後はentitiesを使うようにする

type SnatForDnat struct {
	ExternalInterface string
	InternalIP        string
	ExternalIP        string
}
