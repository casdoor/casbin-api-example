# Casbin API Example for Casdoor

This is a demo project that demonstrates how to use the Exposed Casbin APIs in Casdoor. This example is written in Go and shows how to interact with Casdoor's permission system using Casbin.

## Overview

[Casdoor](https://casdoor.org/) is a UI-first Identity Access Management (IAM) / Single-Sign-On (SSO) platform that supports OAuth 2.0, OIDC, SAML and CAS. It also provides powerful permission management through [Casbin](https://casbin.org/), an authorization library that supports various access control models.

This example demonstrates:
- **Enforce**: Check if a user has permission to perform an action on a resource
- **BatchEnforce**: Check multiple permissions in a single request
- **GetAllObjects**: Get all objects a user can access
- **GetAllActions**: Get all actions a user can perform on an object
- **GetAllRoles**: Get all roles assigned to a user

## Prerequisites

- Go 1.16 or higher
- A running Casdoor instance (or access to a Casdoor server)
- Casdoor credentials (Client ID and Client Secret)

## Installation

1. Clone this repository:
```bash
git clone https://github.com/casdoor/casbin-api-example.git
cd casbin-api-example
```

2. Install dependencies:
```bash
go mod tidy
```

## Configuration

The example can be configured using environment variables:

- `CASDOOR_ENDPOINT`: The URL of your Casdoor instance (default: `http://localhost:8000`)
- `CASDOOR_CLIENT_ID`: Your Casdoor application client ID
- `CASDOOR_CLIENT_SECRET`: Your Casdoor application client secret
- `CASDOOR_ORGANIZATION`: Your Casdoor organization name (default: `built-in`)
- `CASDOOR_APPLICATION`: Your Casdoor application name (default: `app-built-in`)

Example:
```bash
export CASDOOR_ENDPOINT="https://your-casdoor-instance.com"
export CASDOOR_CLIENT_ID="your-client-id"
export CASDOOR_CLIENT_SECRET="your-client-secret"
export CASDOOR_ORGANIZATION="your-org"
export CASDOOR_APPLICATION="your-app"
```

## Running the Example

Run the demo:
```bash
go run main.go
```

## Understanding Casbin

Casbin is an authorization library that supports multiple access control models. In this example, we use the RBAC (Role-Based Access Control) model.

### Model Configuration

The `model.conf` file defines the access control model:
- **request_definition**: Defines the request format (subject, object, action)
- **policy_definition**: Defines the policy format
- **role_definition**: Defines role inheritance
- **policy_effect**: Defines how policies are evaluated
- **matchers**: Defines the matching rules

### Policy Configuration

The `policy.csv` file contains the policies and role assignments:
- `p` lines define permissions (subject, object, action)
- `g` lines define role assignments (user, role)

Example policies:
```
p, alice, data1, read    # alice can read data1
p, alice, data1, write   # alice can write to data1
g, alice, admin          # alice has the admin role
```

## API Examples

### 1. Enforce - Check Permission

Check if a user has permission to perform an action:

```go
allowed, err := client.Enforce(permissionID, []string{"alice", "data1", "read"})
// Returns: true if alice can read data1
```

### 2. BatchEnforce - Check Multiple Permissions

Check multiple permissions in a single request for better performance:

```go
requests := [][]string{
    {"alice", "data1", "read"},
    {"alice", "data1", "write"},
    {"bob", "data2", "read"},
}
results, err := client.BatchEnforce(permissionID, requests)
// Returns: []bool with results for each request
```

### 3. GetAllObjects - Get Accessible Objects

Get all objects a user can access:

```go
objects, err := client.GetAllObjects(permissionID, "alice")
// Returns: []string{"data1", "data2", ...}
```

### 4. GetAllActions - Get Available Actions

Get all actions a user can perform on an object:

```go
actions, err := client.GetAllActions(permissionID, "alice", "data1")
// Returns: []string{"read", "write", ...}
```

### 5. GetAllRoles - Get User Roles

Get all roles assigned to a user:

```go
roles, err := client.GetAllRoles(permissionID, "alice")
// Returns: []string{"admin", ...}
```

## Project Structure

```
casbin-api-example/
├── main.go          # Main demo application
├── model.conf       # Casbin model configuration
├── policy.csv       # Casbin policy rules
├── go.mod           # Go module file
└── README.md        # This file
```

## Additional Resources

- [Casdoor Documentation](https://casdoor.org/docs/)
- [Casdoor Exposed Casbin APIs](https://casdoor.org/docs/permission/exposed-casbin-apis/)
- [Casbin Documentation](https://casbin.org/)
- [Casbin RBAC Model](https://casbin.org/docs/rbac)

## License

This project is licensed under the Apache License 2.0 - see the LICENSE file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.