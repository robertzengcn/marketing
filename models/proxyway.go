package models

type Proxyway interface {
    Proxylist() ([]Proxy,error)
    Createproxy() (error)
    Updateproxy() (error)
    Replaceproxy(string)(error)
}