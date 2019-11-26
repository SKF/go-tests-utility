# go-tests-utility
Go testing related utilities. 

## Packages

### http
### json
### auth
``` go
SignIn(stage, username, password string) (tokens Tokens, err error)
```
### companies
``` go
Create(identityToken, stage, parentNodeID, label, description string) (companyID string, err error)
Delete(identityToken, stage, companyID string) (err error)
```
### disposable-emails
``` go
NewEmailAddress() (emailAddress string, err error)
PollForMessageWithSubject(emailAddress, subject string, fromTimestamp time.Time) (msgAsHTML string, err error)
```
### users
``` go
Create(accessToken, stage, companyID, email string) (createdUser User, password string, err error)
Delete(accessToken, stage, userID string) error
AddUserAccess(identityToken, stage, userID, companyID string) (err error)
AddUserRole(identityToken, stage, userID, role string) (err error)
```