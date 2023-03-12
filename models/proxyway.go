package models

type Proxyway interface {
    Proxylist() ([]Proxy,error)
}