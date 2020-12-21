# env

A golang port of the Ruby dotenv project (which loads env vars from a .env file), plus adapter support for fetching secrets from other resources such as AWS secrets manager

- [Installation](#installation)
- [Usage](#usage)
  - [Load](#Load)
  - [MustLoad](#mustload)
  - [ApplyAdapter](#applyadapter)
  - [LoadSecrets](#loadsecrets)
  - [MustLoadSecrets](#mustloadsecrets)
  - [NewMap](#newmap)
- [Contributing](#contributing)

## Installation

To install the package, you need to install Go and set your Go workspace first.

1. The first need [Go](https://golang.org/) installed (**version 1.13+**), then you can use the Go command below.

```sh
$ go get -u github.com/andreGarvin/env
```

Then you can start using it

## Usage

Create your .env file with all your application environment variables

```.env
APP_NAME=my-cool-app
SECRET=hello world
```

Then in your Go app you can write the following code

```go
package main

import (
  "log"
  "github.com/andreGarvin/env"
)

func main() {
  err := env.Load(".env")
  if err != nil {
    log.Fatal(err)
  }

  fmt.Println(os.Getenv("APP_NAME"))
  fmt.Println(os.Getenv("SECRET"))
}
```

### Load

You can also load more then one .env file name or file path

```golang
err := env.Load(".env", "vault", "another-file-name", "../some/other/file/path")
```

After writing the code above code you can then run the command below and test it

```v
# run example.go
$ go run example.go
```

Here some features I have created for myself based on past challenges I faced when dealing withy env var loading

### MustLoad

Now lets say you want make sure your env vars are loaded on your application and they are not empty, this where `MustLoad` and `RequiredKeys` will be helpful solving this problem.

```env
# Your env file with a missing env var MESSAGE and one that is

APP_NAME=my-cool-app
SECRET=
```

```golang
err := env.MustLoad(".env")
if err != nil {
  log.Fatal(err)
}
```

```
$ go run main.go
2020/12/20 02:40:48 Required keys missing or empty: [SECRET, MESSAGE]
exit status 1
```

### RequiredKeys

Make certain env vars are required in your application use RequiredKeys combined with MustLoad to ensure those env vars are there

```golang
// createda slice of all the env vars that are required for your app
env.RequiredKeys([]string{
  "APP_NAME",
  "MESSAGE",
})
```

### ApplyAdapter

Lets say you want to load some secrets from some secrets manager into your local dev environment for testing or something along those lines. You can use adapters, which is basically code that is ran to pull fetches your secrets and set them into your environment.

```golang
func main() {
  // create a adpater using env.Adapter struct
  adapter := &env.Adapter{
    Pull: foo,
  },

  // adding the adapters
  env.ApplyAdapter(adapter)

  // secrets returned from a adapter will overwrite any env vars declared from the .env file
  err := env.Load(".env")
  if err != nil {
    log.Fatal(err)
  }

  fmt.Println(os.Getenv("MESSAGE"))
}

func foo() (*env.Map, error) {
  e := env.NewMap()

  e.Set("MESSAGE", "FROM SECRET STORE")

  return e, nil
}
```

### LoadSecrets

Now lets say you just want a way for you to load secrets from some secret store into your application in production, well `LoadSecrets` has you covered.

```golang
// only loads secrets from adapters
err := env.LoadSecrets()
if err != nil {
  log.Fatal(err)
}

fmt.Println(os.Getenv("MESSAGE"))
```

### NewMap

This is used to stored env vars before setting them into the environment and to easily join two different maps together

```golang
map1 := env.NewMap()

map1.Set("MAP_1_VAR", "1")

map2 := env.NewMap()
map2.Set("MAP_2_VAR", "2")

map1.SetMap(map2)

fmt.Println(map1.Map)
```

## Contributing

Feel free to send make issues and pull request for any ideas you want to add or making this package even better for developer experience.

1. Fork it
2. Create a branch (`git checkout -b my-changes`)
3. Commit your changes (`git commit -am 'Added some new feature'`)
4. Push to the branch (`git push my-changes`)
5. Create new Pull Request along with a summary of why you added a feature or the changes you made
