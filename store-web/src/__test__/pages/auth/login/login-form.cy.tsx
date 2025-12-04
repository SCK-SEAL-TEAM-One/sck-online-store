import LoginForm from '@/app/auth/login/components/login-form'

describe('<LoginForm />', () => {
  beforeEach(() => {
    // Arrange
    cy.mount(<LoginForm />)
  })

  it('Should show "Password is required." and "Username is required." if fill nothing', () => {
    // Act
    cy.get('#login-btn').click()

    // Assert
    cy.get('#login-username-input-error-txt')
      .should('be.visible')
      .should('have.text', 'Username is required.')
    cy.get('#login-password-input-error-txt')
      .should('be.visible')
      .should('have.text', 'Password is required.')
  })

  it('Should show "Password is required." if fill only username', () => {
    // Act
    cy.get('#login-username-input').type('nattapon.s')
    cy.get('#login-btn').click()

    // Assert
    cy.get('#login-username-input-error-txt').should('not.exist')
    cy.get('#login-password-input-error-txt')
      .should('be.visible')
      .should('have.text', 'Password is required.')
  })

  it('Should show "Username is required." if fill only password', () => {
    // Act
    cy.get('#login-password-input').type('Natta@2025')
    cy.get('#login-btn').click()

    // Assert
    cy.get('#login-username-input-error-txt')
      .should('be.visible')
      .should('have.text', 'Username is required.')
    cy.get('#login-password-input-error-txt').should('not.exist')
  })

  it('Should calls the API when inputs are valid', () => {
    cy.intercept('POST', '**/api/v1/login').as('loginCall')

    // Act
    cy.get('#login-username-input').type('nattapon.s')
    cy.get('#login-password-input').type('Natta@2025')
    cy.get('#login-btn').click()

    // Assert
    cy.wait('@loginCall')
      .its('request.body')
      .should('deep.equal', { username: 'nattapon.s', password: 'Natta@2025' })
  })

  it('Should NOT call the API when inputs are invalid', () => {
    // Arrange
    cy.intercept('POST', '/api/v1/auth/login').as('loginCall')

    // Act
    cy.get('#login-btn').click()

    // Assert
    cy.get('@loginCall.all').should('have.length', 0)
    cy.get('#login-username-input-error-txt')
      .should('be.visible')
      .should('have.text', 'Username is required.')
    cy.get('#login-password-input-error-txt')
      .should('be.visible')
      .should('have.text', 'Password is required.')
  })

  it('Should be able to use "Enter" press key instead of click button', () => {
    cy.intercept('POST', '**/api/v1/login').as('loginCall')

    // Act
    cy.get('#login-username-input').type('nattapon.s')
    cy.get('#login-password-input').type('Natta@2025{enter}')

    // Assert
    cy.wait('@loginCall')
      .its('request.body')
      .should('deep.equal', { username: 'nattapon.s', password: 'Natta@2025' })
  })
})
