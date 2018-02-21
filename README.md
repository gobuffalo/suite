# suite

Suite is a package meant to make testing [gobuffalo.io](http://gobuffalo.io) applications easier.

## Setup

This is the entry point into your unit testing suite. The `Test_ActionSuite(t *testing.T)` function is
compatible with the `go test` command, and it should:

- Create and configure your new test suite instance (`ActionSuite` in this case)
- Call `suite.Run` with the `*testing.T` passed by the Go testing system, and your new `ActionSuite` instance

```go
package actions_test

import (
	"testing"

	"github.com/gobuffalo/suite"
	"github.com/gobuffalo/toodo/actions"
)

type ActionSuite struct {
	*suite.Action
}

func Test_ActionSuite(t *testing.T) {
	as := &ActionSuite{suite.NewAction(actions.App())}
	suite.Run(t, as)
}
```

## Usage

This is where you write your actual test logic. The rules for test names are similar, but not the same, as with `go test`:

- Each test is a method on your `*ActionSuite`
- Test method names should start with `Test` (note the upper case `T`)
- Test methods should have no arguments

A few additional notes:

- To avoid race conditions on the testing database, always use the `ActionSuite` variable called `DB` to access the database (not your production app's database)
- You can access the raw `*testing.T` value if needed with `as.T()`
- `ActionSuite` has support for [`testify`](https://github.com/stretchr/testify)'s [`require` package](https://godoc.org/github.com/stretchr/testify/require) and [`assert` package](https://godoc.org/github.com/stretchr/testify/assert)
- ... So try to use one of those instead packages of using the raw methods on the `*testing.T`
- The default database that `suite` will connect to is called `testing` in your [database.yml](https://github.com/markbates/pop#connecting-to-databases)

```go
package actions_test

import (
	"fmt"

	"github.com/gobuffalo/toodo/models"
)

func (as *ActionSuite) Test_TodosResource_List() {
	todos := models.Todos{
		{Title: "buy milk"},
		{Title: "read a good book"},
	}
	for _, t := range todos {
		err := as.DB.Create(&t)
		as.NoError(err)
	}

	res := as.HTML("/todos").Get()
	body := res.Body.String()
	for _, t := range todos {
		as.Contains(body, fmt.Sprintf("<h2>%s</h2>", t.Title))
	}
}

func (as *ActionSuite) Test_TodosResource_New() {
	res := as.HTML("/todos/new").Get()
	as.Contains(res.Body.String(), "<h1>New Todo</h1>")
}

func (as *ActionSuite) Test_TodosResource_Create() {
	todo := &models.Todo{Title: "Learn Go"}
	res := as.HTML("/todos").Post(todo)
	as.Equal(301, res.Code)
	as.Equal("/todos", res.Location())

	err := as.DB.First(todo)
	as.NoError(err)
	as.NotZero(todo.ID)
	as.NotZero(todo.CreatedAt)
	as.Equal("Learn Go", todo.Title)
}

func (as *ActionSuite) Test_TodosResource_Create_Errors() {
	todo := &models.Todo{}
	res := as.HTML("/todos").Post(todo)
	as.Equal(422, res.Code)
	as.Contains(res.Body.String(), "Title can not be blank.")

	c, err := as.DB.Count(todo)
	as.NoError(err)
	as.Equal(0, c)
}

func (as *ActionSuite) Test_TodosResource_Update() {
	todo := &models.Todo{Title: "Lern Go"}
	verrs, err := as.DB.ValidateAndCreate(todo)
	as.NoError(err)
	as.False(verrs.HasAny())

	res := as.HTML("/todos/%s", todo.ID).Put(&models.Todo{ID: todo.ID, Title: "Learn Go"})
	as.Equal(200, res.Code)

	err = as.DB.Reload(todo)
	as.NoError(err)
	as.Equal("Learn Go", todo.Title)
}
```
