/// <reference types="cypress" />

describe('MSISDNApp tests', () => {
  it('Registration Test', () => {
    cy.visit('http://localhost:8080')
    cy.get('#register').click()
    cy.url().should('include','/register')
    cy.get('#usernameinput').type('fakemail@goodmail.com')
    cy.get('#passwordinput').type('123456Aa!')
    cy.get('#registersubmit').click()
    cy.url().should('not.contain','/register')
  })

  it('Login Test', () => {

    cy.visit('http://localhost:8080/login')
    cy.url().should('include','/login')
    cy.get('#usernameinput').type('fakemail@goodmail.com')
    cy.get('#passwordinput').type('123456Aa!')
    cy.get('#loginsubmit').click()
    cy.url().should('include','/')
  })

  it('Good Lookup test', () => {
    cy.visit('http://localhost:8080/service/lookup')
    cy.url().should('include','/login')
    cy.get('#usernameinput').type('fakemail@goodmail.com')
    cy.get('#passwordinput').type('123456Aa!')
    cy.get('#loginsubmit').click()
    cy.url().should('include','/')
    cy.get('#lookup').click()
    cy.url().should('include','/service/lookup')
    cy.get('#msisdn-input').type('38977123456')
    cy.get('#msisdn-input-submit').click()
    cy.get('#result-wrapper > p').should(($results) => {
      expect($results).to.have.length(4)
      expect($results.eq(0)).to.contain('MNO: A1')
      expect($results.eq(1)).to.contain('Country Code: 389')
      expect($results.eq(2)).to.contain('Subscriber Number: 123456')
      expect($results.eq(3)).to.contain('Country Identifier: mk')
    })
  })

  it('Bad lookup test', () => {
    cy.visit('http://localhost:8080/service/lookup')
    cy.url().should('include','/login')
    cy.get('#usernameinput').type('fakemail@goodmail.com')
    cy.get('#passwordinput').type('123456Aa!')
    cy.get('#loginsubmit').click()
    cy.url().should('include','/')
    cy.get('#lookup').click()
    cy.url().should('include','/service/lookup')
    cy.get('#msisdn-input').type('123456')
    cy.get('#msisdn-input-submit').click()
    cy.get('#result-wrapper > p').should(($results) => {
      expect($results).to.have.length(1)
      expect($results.eq(0)).to.contain('Error')
    })
  })

})