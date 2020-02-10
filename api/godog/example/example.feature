Feature: example feature
  show how to use the go test utility for creating go dog features

  @example
  Scenario: create user and execute request
    Given the companies: "myCompany"
    And the users: "foo" for company: ".companies.myCompany.id"

    When the scenario creates a "GET :: /users/me" request
    And  sets request header parameter "Authorization" to ".users.foo.tokens.accessToken"
    And  executes the request

    Then the response code should be: 200, "OK"
    And the response body value ".data.id" equals ".user.foo.id"
    And the response body value ".data.companyId" equals ".companies.myCompany.id"
    And the response body value ".data.email" equals ".users.foo.username"
    And the response body value ".data.givenName" equals ".users.foo.givenName"
    And the response body value ".data.surname" equals ".users.foo.surname"
    And the response body value ".data.language" equals ".users.foo.language"
    And the response body value ".data.status" equals ".users.foo.status"
