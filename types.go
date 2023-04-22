package main

import "go.mau.fi/whatsmeow"

type Validator func(interface{}, *whatsmeow.Client) (bool, error)
type Executor func(interface{}, *whatsmeow.Client) []error
