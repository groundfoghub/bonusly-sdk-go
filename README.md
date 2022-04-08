--- 
⚠️ Important: This Go package is in very early development. Do not use it in production. Everything about this package might change in the future.

---

[![Go Reference](https://pkg.go.dev/badge/github.com/groundfoghub/bonusly-sdk-go.svg)](https://pkg.go.dev/github.com/groundfoghub/bonusly-sdk-go)

# (Unofficial) Bonus.ly SDK for Go

At [Groundfog](https://groundfog.cloud) we value employee recognition, which is why we love using [Bonus.ly](https://bonus.ly). To help us automate a few things we wanted to create a Go package that we can use to create reports or give our employees automated bonuses. 

Since we also love sharing with the community and giving back, we decided to build our SDK in the open and allow others to also build something awesome with it as well.

## Installation

```shell
go get github.com/groundfoghub/bonusly-sdk-go
```

## Updating

```shell
go get -u github.com/groundfoghub/bonusly-sdk-go
```

## Getting Started

Before you start with the examples, you need to get a valid access token. To get a token login to Bonus.ly and go to https://bonus.ly/api to create a new access token.

**Get all users**

```go
config := bonusly.Configuration{Token: "<your-access-token>"}
client := bonusly.New(config)

var users []bonusly.User

paginator := bonusly.NewListUsersPaginator(client, nil)
for paginator.HasMorePages() {
    output, err := paginator.NextPage(context.TODO())
    if err != nil {
        return
    }

    users = append(users, output.Users...)
}

fmt.Printf("Found %d users\n", len(users))
```

**Create a bonus**

You need a token that allows write access for this example to work.

```go
config := bonusly.Configuration{Token: "<your-access-token>"}
client := bonusly.New(config)

params := bonusly.CreateBonusInput{
    GiverEmail: "leia@examplecorp.com",
    Receivers:  []string{"luke@examplecorp.com"},
    Reason:     "For destroying the Death Star",
    Amount:     25,
}

_, err := client.CreateBonus(context.TODO(), &params)
if err != nil {
    fmt.Println("create bonus: ", err)
}
```