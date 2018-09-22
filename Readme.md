# go-rest-api
Simple examples of HTTP - JSON REST API in Go Language(golang) using PostgreSQL database.
- USER MANAGEMENT
- TODOLIST MANAGEMENT
---

USER MANAGEMENT
---
It is designed and developed to implement following operations:
- Create User: To create a new user
- Get User: To get the required user using its UserID
- GetAll User: To get all the users from database
- Edit User: To edit the user details
- Delete User: To delete a user

A user has following information attributes:
- First Name
- Last Name
- Date of Birth (Optional)
- Email
- Phone Number

---

TODOLIST MANAGEMENT
---
It is designed and developed to implement following operations:
- Add Todo List: To add a new todo list
- Delete Todo List: To delete an already present todo list
- Edit Todo List Name: To update the name of an already present todo list
- Add Todo Item: To add an item in a todo list
- Delete Todo List Item: To delete an item of a todo list
- Get Todo List Item: To get an item of a todo list
- Update Todo Item: To update an item of a list
- Get Todo List : To get the whole todo list

A TodoItem has following information attributes:
- ID: Item ID
- Value: Item Value/Description
- Completed: Item Status

A TodoList has following information attributes:
- ID: List ID
- Items: List of TodoItems in the TodoList
- Name: List Name/Description

---

Both the API services are protected using [Basic Auth](https://en.wikipedia.org/wiki/Basic_access_authentication) with following credentials: 
- Username: mavis
- Password: shivam

# Contributing
Changes and improvements are more than welcome! 
Feel free to fork and open a pull request. 
And Please make your changes in a specific branch and request to pull into master! If you can, please make sure the game fully works before sending the PR, as that will help speed up the process.

# License
GO-REST-API is licensed under the [MIT license](https://github.com/Shivam010/go-rest-api/blob/master/LICENSE).