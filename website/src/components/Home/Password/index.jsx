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

  const type = 'pass'
  const secretID = get(secret, 'id')
  const [note, setNote] = useState(get(secret, 'note', ''))
  const [title, setTitle] = useState(get(secret, 'title', ''))
  const [login, setLogin] = useState(get(secret, 'login', ''))
  const [password, setPassword] = useState(get(secret, 'password', ''))

  const setValue = (event) => {
    const field = get(event, 'target.id', '')
    const value = get(event, 'target.value', '')

    const changeValueHandler = {
      note     : setNote,
      title    : setTitle,
      login    : setLogin,
      password : setPassword,
    }[field] || noop

    changeValueHandler(value)
  }

  const onSubmitButtonClick = async () => {
    const apiHandler = isUndefined(secretID) ? createSecret : updateSecret

    try {
      await apiHandler({
        id: secretID,
        type,
        title,
        note,
        login,
        password,
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
            <label htmlFor="login" className="form-label">Login</label>
            <input type="text" className="form-control" id="login" value={login} onChange={setValue} />
            <small className="text-body-secondary">Full name as displayed on card</small>
          </div>

          <div className="col-md-6">
            <label htmlFor="password" className="form-label">Password</label>
            <input type="text" className="form-control" id="password" value={password} onChange={setValue} />
          </div>

          <div className="col-md-12">
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
