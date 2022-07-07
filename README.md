
:warning: **Important**: This Go package is in very early development. Do not use it in production. Everything about this package might change in the future.

---

[![Go Reference](https://pkg.go.dev/badge/github.com/groundfoghub/bonusly-sdk-go.svg)](https://pkg.go.dev/github.com/groundfoghub/bonusly-sdk-go)

# (Unofficial) Bonus.ly SDK for Go

At [Groundfog](https://groundfog.cloud) we value employee recognition, which is why we love using [Bonus.ly](https://bonus.ly). To help us automate a few things we wanted to create a Go package that we can use to create reports or give our employees automated bonuses. 

Since we also love sharing with the community and giving back, we decided to build our SDK in the open and allow others to also build something awesome with it as well.

## Table of Contents

- [Installation](#package-installation)
- [Getting Started](#sparkles-getting-started)
- [Implementation Status](#white_check_mark-implementation-status)
- [Documentation](#pencil2-documentation)
- [Contributing](#sparkling_heart-contributing)

## :package: Installation
[(Back to top)](#table-of-contents)

```shell
go get github.com/groundfoghub/bonusly-sdk-go
```

To update the SDK run the following command:

```shell
go get -u github.com/groundfoghub/bonusly-sdk-go
```

## :sparkles: Getting Started
[(Back to top)](#table-of-contents)

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
**Receive Achievements**

You'll need an admin token to make this example work

```go
config := bonusly.Configuration{Token: "<your-access-token>"}
client := bonusly.New(config)

// Receives 50 most recent achievements 
achievements, err := bonusly.GetAchivements()

if err != nil {
    fmt.Println("create bonus: ", err)
}

for _, achievement := range achievements {
    // Do something with the achievement
}

```

## :white_check_mark: Implementation Status
[(Back to top)](#table-of-contents)

The SDK does not yet cover all of the Bonus.ly API yet. Based on the [official API documentation](https://bonusly.docs.apiary.io/#), the following list provides an overview of all the features and their current implementation status.

At the moment we are working on implementing all the features except SCIM, which will follow later. Current progress can be tracked using the [All Non-SCIM Features Implemented](https://github.com/groundfoghub/bonusly-sdk-go/milestone/1) milestone.

Legend:
* :white_check_mark: Implemented.
* :warning: Partially implemented or can not be exactly fully implemented.
* :no_entry: Not implemented yet.

**Achievements**
* :white_check_mark: List Achievements

**Analytics**
* :no_entry: Trends | Index (Issue: [#7](https://github.com/groundfoghub/bonusly-sdk-go/issues/7))
* :no_entry: Leaderboards | Index (Issue: [#8](https://github.com/groundfoghub/bonusly-sdk-go/issues/8))

**API Keys**
* :no_entry: List API Keys (Issue: [#9](https://github.com/groundfoghub/bonusly-sdk-go/issues/9))
* :no_entry: Create API Key (Issue: [#10](https://github.com/groundfoghub/bonusly-sdk-go/issues/10))
* :no_entry: Cancel API Key (Issue: [#11](https://github.com/groundfoghub/bonusly-sdk-go/issues/11))

**Bonuses**
* :no_entry: List Bonuses (Issue: [#12](https://github.com/groundfoghub/bonusly-sdk-go/issues/12))
* :white_check_mark: Create a Bonus
* :warning: Create a Bonus with separate fields fo reason, hashtag, receiver and amount
* :no_entry: Retrieve a Bonus (Issue: [#13](https://github.com/groundfoghub/bonusly-sdk-go/issues/13))
* :no_entry: Update a Bonus (Issue: [#14](https://github.com/groundfoghub/bonusly-sdk-go/issues/14))
* :no_entry: Delete a Bonus (Issue: [#15](https://github.com/groundfoghub/bonusly-sdk-go/issues/15))

**Company**
* :no_entry: Retrieve a Company (Issue: [#16](https://github.com/groundfoghub/bonusly-sdk-go/issues/16))
* :no_entry: [ADMIN] Update a Company (Issue: [#17](https://github.com/groundfoghub/bonusly-sdk-go/issues/17))

**Redemptions**
* :white_check_mark: List Redemptions
* :white_check_mark: Retrieve a Redemption

**Rewards**
* :white_check_mark: List Rewards
* :white_check_mark: Retrieve a Reward

**SCIM**
* :no_entry: List users
* :no_entry: Retrieve a user
* :no_entry: Create a user
* :no_entry: Update an existing user
* :no_entry: Activate or deactivate a user
* :no_entry: Get metadata about the Bonusly SCIM API
* :no_entry: List the SCIM resource types supported by Bonusly
* :no_entry: List the SCIM schemas supported by Bonusly

**Users**
* :white_check_mark: List Users
* :white_check_mark: Retrieve a User
* :no_entry: Me (Issue: [#18](https://github.com/groundfoghub/bonusly-sdk-go/issues/18))
* :no_entry: Autocomplete (Issue: [#19](https://github.com/groundfoghub/bonusly-sdk-go/issues/19))
* :no_entry: Bonuses (Issue: [#20](https://github.com/groundfoghub/bonusly-sdk-go/issues/20))
* :no_entry: Achievements (Issue: [#21](https://github.com/groundfoghub/bonusly-sdk-go/issues/21))
* :no_entry: Redemptions (Issue: [#22](https://github.com/groundfoghub/bonusly-sdk-go/issues/22))
* :no_entry: Create a Redemption (Issue: [#23](https://github.com/groundfoghub/bonusly-sdk-go/issues/23))
* :no_entry: [ADMIN] Create a User (Issue: [#24](https://github.com/groundfoghub/bonusly-sdk-go/issues/24))
* :no_entry: [ADMIN] Update a User (Issue: [#25](https://github.com/groundfoghub/bonusly-sdk-go/issues/25))
* :no_entry: [ADMIN] Deactivate a User(Issue: [#26](https://github.com/groundfoghub/bonusly-sdk-go/issues/26))

**Webhooks**
* :white_check_mark: List Webhooks
* :white_check_mark: Create Webhook
* :white_check_mark: Update Webhook
* :white_check_mark: Remove Webhook

## :pencil2: Documentation
[(Back to top)](#table-of-contents)

The official Go package documentation can be found at [pkg.go.dev](https://pkg.go.dev/github.com/groundfoghub/bonusly-sdk-go). 

## :sparkling_heart: Contributing
[(Back to top)](#table-of-contents)

If you found a bug, have a feature suggestion or just want to help us build the SDK, feel free to [file an issue](https://github.com/groundfoghub/bonusly-sdk-go/issues/new) or create a pull requests. Contributions are welcome.
