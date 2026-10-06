# Casdoor Casbin API Example

[![Build](https://github.com/casdoor/casbin-api-example/actions/workflows/build.yml/badge.svg)](https://github.com/casdoor/casbin-api-example/actions/workflows/build.yml)
[![License](https://img.shields.io/github/license/casdoor/casbin-api-example)](https://github.com/casdoor/casbin-api-example/blob/master/LICENSE)
[![Discord](https://img.shields.io/discord/1022748306096537660?logo=discord&label=discord&color=5865F2)](https://discord.gg/5rPsrAzK7S)

A Go command-line example of the [exposed Casbin APIs](https://casdoor.ai/docs/permission/exposed-casbin-apis) of [Casdoor](https://casdoor.ai/): your application asks Casdoor whether a user is allowed to do something, and Casdoor decides with the permissions you manage in it.

The example uses [casdoor-go-sdk](https://github.com/casdoor/casdoor-go-sdk) and shows:

| API              | What it does                                                    |
|------------------|-----------------------------------------------------------------|
| `enforce`        | Checks whether a user can do an action on a resource            |
| `batch-enforce`  | Checks several requests in one call                             |
| `get-all-objects`| Lists the resources in the permissions that apply to a user     |
| `get-all-actions`| Lists the actions in the permissions that apply to a user       |
| `get-all-roles`  | Lists the roles of a user                                       |

## What the example does

1. Creates a user `alice`, a model (from [model.conf](model.conf)), a role that `alice` belongs to, and a permission that allows the role to `read` and `write` `data1`.
2. Calls the five APIs above for `alice`, and for `bob`, who has no role.
3. Deletes everything it created.

```
Setup: alice has the role, the role is allowed to read and write data1

1. Enforce
   [built-in/casbin-api-example-alice data1 read] -> true
   [built-in/casbin-api-example-alice data2 read] -> false
   [built-in/casbin-api-example-bob data1 read] -> false

2. BatchEnforce
   [built-in/casbin-api-example-alice data1 read] -> true
   [built-in/casbin-api-example-alice data1 write] -> true
   [built-in/casbin-api-example-alice data1 delete] -> false
   [built-in/casbin-api-example-bob data1 read] -> false

3. GetAllObjects: [data1]
4. GetAllActions: [read write]
5. GetAllRoles:   [casbin-api-example-role]
```

## Prerequisites

- Go 1.23+
- A Casdoor server, see [Casdoor installation](https://casdoor.ai/docs/basic/server-installation)
- The client ID and client secret of a Casdoor application. The example creates and deletes objects, so the application must belong to the `built-in` organization (which makes it a global admin) or to the organization you run the example in.

## Configuration

The example is configured with environment variables:

| Variable                | Default                 | Description                                     |
|-------------------------|-------------------------|-------------------------------------------------|
| `CASDOOR_ENDPOINT`      | `http://localhost:8000` | Casdoor server URL                              |
| `CASDOOR_CLIENT_ID`     |                         | Client ID of the application                    |
| `CASDOOR_CLIENT_SECRET` |                         | Client secret of the application                |
| `CASDOOR_ORGANIZATION`  | `built-in`              | Organization the example creates its objects in |
| `CASDOOR_APPLICATION`   | `app-built-in`          | Name of the application                         |

## Run

```shell
git clone https://github.com/casdoor/casbin-api-example
cd casbin-api-example

export CASDOOR_ENDPOINT="http://localhost:8000"
export CASDOOR_CLIENT_ID="<client-id>"
export CASDOOR_CLIENT_SECRET="<client-secret>"

go run .
```

## How it works

The model in [model.conf](model.conf) is a Casbin RBAC model: a request is `subject, object, action`, and it's allowed when the subject, or one of its roles, has a policy for the object and the action.

A Casdoor permission generates the policies: its users and roles are the subjects, its resources are the objects and its actions are the actions. So the permission of the example becomes the policies `role, data1, read` and `role, data1, write`, and the role's users become `alice, role`.

In your code, call `Enforce()` with the ID of the permission and the request:

```go
client := casdoorsdk.NewClient(endpoint, clientId, clientSecret, certificate, organization, application)

allowed, err := client.Enforce("built-in/my-permission", "", "", "", "",
	casdoorsdk.CasbinRequest{"built-in/alice", "data1", "read"})
```

Instead of one permission, `Enforce()` can also check all the permissions of a model, a resource, an enforcer or an owner, see the [documentation](https://casdoor.ai/docs/permission/exposed-casbin-apis).

## Resources

- [Casdoor documentation](https://casdoor.ai/docs/overview)
- [Exposed Casbin APIs](https://casdoor.ai/docs/permission/exposed-casbin-apis)
- [casdoor-go-sdk](https://github.com/casdoor/casdoor-go-sdk)
- [Casbin documentation](https://casbin.org/docs/overview)

## License

[Apache-2.0](LICENSE)
