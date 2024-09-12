package server

type Server interface {
	Run(address string) error
}
