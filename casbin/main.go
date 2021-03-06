package main

import (
	"fmt"

	"github.com/casbin/casbin"
)

func main() {
	e := casbin.NewEnforcer("model.conf", "policy.csv")

	test(e, "luoji", "/a", "write")
	test(e, "luoji", "/member/a", "write")
	test(e, "tom", "/a", "write")
	test(e, "tom", "/member/a", "write")

	//test(e, "anonymous", "/register", "read")
}

func test(e *casbin.Enforcer, sub, obj, act string) {
	//sub := "luoji" // the user that wants to access a resource.
	//obj := "/abc"  // the resource that is going to be accessed.
	//act := "read"  // the operation that the user performs on the resource.

	ret := e.Enforce(sub, obj, act)
	fmt.Println("-->", ret)
	/*
		if DeepEnforce(e, sub, obj, act) == true {
			// permit alice to read data1
			fmt.Printf("%s %s %s ok.\n", sub, act, obj)
		} else {
			fmt.Printf("%s %s %s no.\n", sub, act, obj)
			// deny the request, show an error
		}
	*/
}

func DeepEnforce(e *casbin.Enforcer, sub, obj, act string) bool {
	if e.Enforce(sub, obj, act) == true {
		return true
	}

	roles := e.GetRolesForUser(sub)
	for _, role := range roles {
		if e.Enforce(role, obj, act) == true {
			return true
		}
	}

	users := e.GetUsersForRole(sub)
	for _, user := range users {
		if e.Enforce(user, obj, act) == true {
			return true
		}
	}

	return false
}
