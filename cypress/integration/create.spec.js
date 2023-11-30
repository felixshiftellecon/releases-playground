describe ('Create snippet', () => {
  const username = 'createForTest'
  const email = 'create@example.com'
  const password = 'validPa$$word'
  const title = 'Ghost'
  const content = 'A ghost from the past\nHaunting present and future\nLegacy systems\n– Dave'

  before(() => {
    cy.signupForTest(username, email, password)
  })

  context('not logged in', () => {
    it('redirects user to signup', () => {
      cy.visit('https://localhost:4000/snippet/create')

      cy.location('pathname').should('eq', '/user/login')
    })
  })

  context('logged in', () => {

    beforeEach(() => {
      cy.loginForTest(email, password)
      cy.visit('https://localhost:4000/snippet/create')
    })

    it('loads', () => {
      cy.location('pathname').should('eq', '/snippet/create')
    })

    it('requires a title', () => {
      cy.get('textarea[name=content]').type(content)
      cy.get('input[type=submit]').click()

      cy.get('.error').should('have.text', 'This field cannot be blank')
      cy.get('textarea[name=content]').should('have.value', content)
      cy.location('pathname').should('eq', '/snippet/create')
    })

    it('requires content', () => {
      cy.get('input[name=title]').type(title)
      cy.get('input[type=submit]').click()

      cy.get('.error').should('have.text', 'This field cannot be blank')
      cy.get('input[name=title]').should('have.value', title)
      cy.location('pathname').should('eq', '/snippet/create') 
    })

    it('publishes a snippet', () => {
      cy.get('input[name=title]').type(title)
      cy.get('textarea[name=content]').type(content)
      cy.get('input[type=submit]').click()

      cy.get('div[class=flash]').should('have.text', 'Snippet successfully created!')
      cy.get('div[class=snippet]').should('contain', title)
      cy.get('div[class=snippet]').should('contain', content)
      cy.location('pathname').should('match', /\/snippet\/[0-9]*$/)
    })
  })
})