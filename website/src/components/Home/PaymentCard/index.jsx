import get from 'lodash/get'
import noop from 'lodash/noop'
import isUndefined from 'lodash/isUndefined'
import {useContext, useState} from 'react'
import ApplicationContext from '../../../context'

import './style.css'

export default () => {
  const {
    secret,
    navigateTo,
    createSecret,
    updateSecret,
    setAlertMessage,
  } = useContext(ApplicationContext)

  const type = 'card'
  const secretID = get(secret, 'id')
  const [note, setNote] = useState(get(secret, 'note', ''))
  const [title, setTitle] = useState(get(secret, 'title', ''))
  const [securityCode, setSecurityCode] = useState(get(secret, 'security_code', ''))
  const [cardNumber, setCardNumber] = useState(get(secret, 'card_number', ''))
  const [expiration, setExpiration] = useState(get(secret, 'expiration', ''))
  const [cardholderName, setCardholderName] = useState(get(secret, 'cardholder_name', ''))

  const setValue = (event) => {
    const field = get(event, 'target.id', '')
    const value = get(event, 'target.value', '')

    const changeValueHandler = {
      note       : setNote,
      title      : setTitle,
      cardnumber : setCardNumber,
      expiration : setExpiration,
      cvv        : setSecurityCode,
      cardholder : setCardholderName,
    }[field] || noop

    changeValueHandler(value)
  }

  const onSubmitButtonClick = async () => {
    // if secret has 'id' fieldm then we're 'updating', otherwise - 'creating'
    const apiHandler = isUndefined(secretID) ? createSecret : updateSecret

    try {
      await apiHandler({
        id              : secretID,
        type,
        title,
        cardholder_name : cardholderName,
        card_number     : cardNumber,
        expiration,
        security_code   : securityCode,
        note,
      }, secretID)

      navigateTo('home')
    } catch (error) {
      console.warn(error.message)
      setAlertMessage(`Error: ${error.message}`)
    }
  }

  return (
    <div className="kms-secret-edit">
      <form className="kms-secret-edit__form">
        <div className="row gy-3">
          <div className="col-md-12">
            <label htmlFor="title" className="form-label">Title</label>
            <input type="text" className="form-control" id="title" value={title} onChange={setValue} />
          </div>

          <div className="col-md-6">
            <label htmlFor="cardholder" className="form-label">Name on card</label>
            <input type="text" className="form-control" id="cardholder" value={cardholderName} onChange={setValue} />
            <small className="text-body-secondary">Full name as displayed on card</small>
          </div>

          <div className="col-md-6">
            <label htmlFor="cardnumber" className="form-label">Credit card number</label>
            <input type="text" className="form-control" id="cardnumber" value={cardNumber} onChange={setValue} />
          </div>

          <div className="col-md-3">
            <label htmlFor="expiration" className="form-label">Expiration</label>
            <input type="date" className="form-control" id="expiration" value={expiration} onChange={setValue} />
          </div>

          <div className="col-md-3">
            <label htmlFor="cvv" className="form-label">CVV</label>
            <input type="text" className="form-control" id="cvv" value={securityCode} onChange={setValue} />
          </div>

          <div className="col-md-6">
            <label htmlFor="note" className="form-label">Add a few notes here</label>
            <textarea className="form-control" id="note" value={note} onChange={setValue} />
          </div>
        </div>
      </form>
      <div className="kms-secret-edit__controls">
        <button type="button" className="btn btn-primary kms-button-add" onClick={onSubmitButtonClick}>
          Save
        </button>
        <button type="button" className="btn btn-danger kms-button-add" onClick={() => navigateTo('home')}>
          Cancel
        </button>
      </div>
    </div>
  )
}
