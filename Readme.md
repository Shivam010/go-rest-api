# go-rest-api
Simple examples of HTTP - JSON REST API in Go Language(golang).

The first example API is based on user-management.

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

The API service is protected using [Basic Auth](https://en.wikipedia.org/wiki/Basic_access_authentication) with following credentials: 
- Username: mavis
- Password: shivam
