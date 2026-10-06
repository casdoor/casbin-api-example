// Copyright 2025 The Casdoor Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
)

const (
	modelName      = "casbin-api-example-model"
	roleName       = "casbin-api-example-role"
	permissionName = "casbin-api-example-permission"
	userName       = "casbin-api-example-alice"
)

func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func check(err error) {
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func getStrings(client *casdoorsdk.Client, action string, userId string) []string {
	bytes, err := client.DoGetBytes(client.GetUrl(action, map[string]string{"userId": userId}))
	check(err)

	var res []string
	err = json.Unmarshal(bytes, &res)
	check(err)
	return res
}

func main() {
	organization := getEnv("CASDOOR_ORGANIZATION", "built-in")
	client := casdoorsdk.NewClient(
		getEnv("CASDOOR_ENDPOINT", "http://localhost:8000"),
		getEnv("CASDOOR_CLIENT_ID", ""),
		getEnv("CASDOOR_CLIENT_SECRET", ""),
		"",
		organization,
		getEnv("CASDOOR_APPLICATION", "app-built-in"),
	)

	modelText, err := os.ReadFile("model.conf")
	check(err)

	alice := organization + "/" + userName
	bob := organization + "/casbin-api-example-bob"
	permissionId := organization + "/" + permissionName

	user := &casdoorsdk.User{Owner: organization, Name: userName, DisplayName: "Alice"}
	model := &casdoorsdk.Model{Owner: organization, Name: modelName, DisplayName: modelName, ModelText: string(modelText)}
	role := &casdoorsdk.Role{Owner: organization, Name: roleName, DisplayName: roleName, Users: []string{alice}, IsEnabled: true}
	permission := &casdoorsdk.Permission{
		Owner:        organization,
		Name:         permissionName,
		DisplayName:  permissionName,
		Roles:        []string{organization + "/" + roleName},
		Model:        organization + "/" + modelName,
		ResourceType: "Custom",
		Resources:    []string{"data1"},
		Actions:      []string{"read", "write"},
		Effect:       "Allow",
		IsEnabled:    true,
		State:        "Approved",
	}

	cleanUp := func() {
		_, _ = client.DeletePermission(permission)
		_, _ = client.DeleteRole(role)
		_, _ = client.DeleteModel(model)
		_, _ = client.DeleteUser(user)
	}
	cleanUp()
	defer cleanUp()

	fmt.Println("Setup: alice has the role, the role is allowed to read and write data1")
	_, err = client.AddUser(user)
	check(err)
	_, err = client.AddModel(model)
	check(err)
	_, err = client.AddRole(role)
	check(err)
	_, err = client.AddPermission(permission)
	check(err)

	fmt.Println()
	fmt.Println("1. Enforce")
	for _, request := range []casdoorsdk.CasbinRequest{
		{alice, "data1", "read"},
		{alice, "data2", "read"},
		{bob, "data1", "read"},
	} {
		allowed, err := client.Enforce(permissionId, "", "", "", "", request)
		check(err)
		fmt.Printf("   %v -> %v\n", request, allowed)
	}

	fmt.Println()
	fmt.Println("2. BatchEnforce")
	requests := []casdoorsdk.CasbinRequest{
		{alice, "data1", "read"},
		{alice, "data1", "write"},
		{alice, "data1", "delete"},
		{bob, "data1", "read"},
	}
	results, err := client.BatchEnforce(permissionId, "", "", "", "", requests)
	check(err)
	for i, request := range requests {
		fmt.Printf("   %v -> %v\n", request, results[0][i])
	}

	fmt.Println()
	fmt.Println("3. GetAllObjects:", getStrings(client, "get-all-objects", alice))
	fmt.Println("4. GetAllActions:", getStrings(client, "get-all-actions", alice))
	fmt.Println("5. GetAllRoles:  ", getStrings(client, "get-all-roles", alice))
}
