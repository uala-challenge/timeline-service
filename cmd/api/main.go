package main

import "github.com/uala-challenge/simple-toolkit/pkg/simplify/app_builder"

func main() {
	builder := NewAppBuilder()
	application := app_builder.Apply(builder)
	err := application.Run()
	if err != nil {
		panic(err)
	}
}
