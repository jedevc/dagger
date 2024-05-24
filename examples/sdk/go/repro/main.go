package main

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
)

func main() {
	ctx := context.Background()

	// this forces dagger connect to create a *new* session subprocess
	os.Unsetenv("DAGGER_SESSION_PORT")

	// this issue *seems* to have something to do with telemetry. unsetting
	// this as well "makes everything work normally again".
	// os.Unsetenv("TRACEPARENT")

	client, err := dagger.Connect(ctx)
	if err != nil {
		panic(err)
	}
	p, err := client.DefaultPlatform(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println(p)
	fmt.Println("trying to close")
	err = client.Close()
	fmt.Println("closed!")
	if err != nil {
		panic(err)
	}
}
