package main

import (
    "context"
    "fmt"
    "os"

    "dagger.io/dagger"
)

func main() {
    if err := build(context.Background()); err != nil {
        fmt.Println(err)
    }
}

// func build(ctx context.Context) error {
//     // initialize Dagger client
//     client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
//     if err != nil {
//         return err
//     }
//     defer client.Close()

//     client.Container().From("alpine@sha256:13b7e62e8df80264dbb747995705a986aa530415763a6c58f84a3ca8af9a5bcd").
//         WithExec([]string{"echo", "hello1"}).Sync(ctx)

//     return nil
// } 

func build(ctx context.Context) error {
    // initialize Dagger client
    client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
    if err != nil {
        return err
    }
    defer client.Close()

    client.Container().
    From("alpine:latest").
        // From("alpine").
        WithDirectory("/root/bazr", client.Host().Directory(".")).
        WithExec([]string{"cat", "/root/bazr/main.go"}).Sync(ctx)

    return nil
}

