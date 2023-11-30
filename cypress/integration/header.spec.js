describe ('The Header', () => {
  const username = 'headerForTest'
  const email = 'header@example.com'
  const password = 'validPa$$word'

  before(() => {
    cy.signupForTest(username, email, password)
  })

  beforeEach(() => {
    cy.visit('https://localhost:4000/')
  })

  it('user not logged in', () => {
    cy.get('a[href="/user/signup"]').should('exist')
    cy.get('a[href="/user/login"]').should('exist')
    cy.get('a[href="/"]').should('exist')

    cy.get('form[action="/user/logout"]').should('not.exist')
    cy.get('a[href="/snippet/create"]').should('not.exist')
  })

  it('user logged in', () => {
    cy.loginForTest(email, password)

    cy.get('a[href="/user/signup"]').should('not.exist')
    cy.get('a[href="/user/login"]').should('not.exist')

    cy.get('a[href="/"]').should('exist')
    cy.get('form[action="/user/logout"]').should('exist')
    cy.get('a[href="/snippet/create"]').should('exist')
  })
})