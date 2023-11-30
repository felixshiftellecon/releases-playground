Cypress.Commands.add("signupForTest", (username, email, password) => {
  cy.visit('https://localhost:4000/user/signup')
  cy.get('input[name=name]').type(username)
  cy.get('input[name=email]').type(email)
  cy.get('input[name=password]').type(password)
  cy.get('input[type=submit]').click()
})

Cypress.Commands.add("loginForTest", (email, password) => {
  cy.visit('https://localhost:4000/user/login')
  cy.get('input[name=email]').type(email)
  cy.get('input[name=password]').type(password)
  cy.get('input[type=submit]').click()
})