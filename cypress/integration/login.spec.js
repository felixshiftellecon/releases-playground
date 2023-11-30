describe ('The Login Page', () => {

  const username = 'loginForTest'
  const email = 'login@example.com'
  const password = 'validPa$$word'

  before(() => {
    cy.signupForTest(username, email, password)
  })

  beforeEach(() => {
    cy.visit('https://localhost:4000/user/login')
  })

  it('loads', () => {
    cy.location('pathname').should('eq', '/user/login')
  })

  it('requires an email address', () => {
    cy.get('input[name=password]').type(`${password}{enter}`)

    cy.get('.error').should('have.text', 'Email or Password is incorrect')
    cy.get('input[name=password]').should('be.empty')
    cy.location('pathname').should('eq', '/user/login')
  })

  it('requires a password', () => {
    cy.get('input[name=email]').type(email)
    cy.get('input[type=submit]').click()

    cy.get('.error').should('have.text', 'Email or Password is incorrect')
    cy.get('input[name=email]').should('have.value', email)
    cy.location('pathname').should('eq', '/user/login')
  })

  it('user can login', () => {

    //  This doesn't work for some reason, even on PostMan
    //   cy.request({
    //   method: 'POST', 
    //   url: 'https://localhost:4000/user/signup', 
    //   followRedirect: false,
    //   form: true,
    //   body: {
    //     "name": `${name}`,
    //     "email": `${email}`,
    //     "password": `${password}`,
    //   },
    // })

    cy.loginForTest(email, password)

    cy.location('pathname').should('eq', '/snippet/create')
  })
})