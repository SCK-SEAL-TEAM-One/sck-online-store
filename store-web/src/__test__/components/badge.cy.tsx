import Badge from '../../components/badge'

describe('<Badge />', () => {
  it('mounts', () => {
    cy.mount(<Badge total={0} />)
  })

  it('id badge-test should display 400', () => {
    cy.mount(<Badge id="badge-test" total={400} />)
    cy.get('#badge-test ').should('have.text', '400')
  })
})
