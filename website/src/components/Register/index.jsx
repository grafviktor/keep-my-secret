import some from 'lodash/some'
import {useContext, useState} from 'react'
import Logo from '../Logo'
import Version from '../Version'
import ApplicationContext from '../../context'
import {isTokenExpired} from '../../utils'
import './style.css'

export default () => {
  const [username, setUsername] = useState('user@localhost')
  const [password, setPassword] = useState('12345')
  const [password2, setPassword2] = useState('12345')
  const {
    api,
    navigateTo,
    setAccessToken,
    setAlertMessage,
  } = useContext(ApplicationContext)

  const onRegisterButtonClick = async (ev) => {
    ev.preventDefault()

    if (!some([username, password])) {
      setAlertMessage('Username or password cannot be empty')

      return
    }

    if (password !== password2) {
      setAlertMessage('The passwords do not match')

      return
    }

    try {
      const {data: response} = await api.register(username, password)

      setAccessToken(response.data)

      if (!isTokenExpired(response.data)) {
        navigateTo('home')
      } else {
        throw Error('there was a problem with logging in, because /register endpoint does not return token')
      }
    } catch (error) {
      console.warn(error)
      setAlertMessage(`Error: ${error.message}`)
    }
  }

  const onSignInClick = (ev) => {
    ev.preventDefault()

    navigateTo('login')
  }

  return (
    <div className="form-register w-100 m-auto">
      <Logo />
      <form>
        <h1 className="h3 mb-3 fw-normal">
          Please register or
          {' '}
          <a href="register" onClick={onSignInClick}>sign in</a>
        </h1>

        <div className="form-floating">
          <input
            type="email"
            className="form-control"
            id="floatingInput"
            placeholder="name@example.com"
            onChange={(event) => { setUsername(event.target.value) }}
            value={username}
          />
          <label htmlFor="floatingInput">Email address</label>
        </div>
        <div className="form-floating">
          <input
            type="password"
            className="form-control"
            id="floatingPassword"
            placeholder="Password"
            onChange={(event) => { setPassword(event.target.value) }}
            value={password}
          />
          <label htmlFor="floatingPassword">Password</label>
        </div>
        <div className="form-floating">
          <input
            type="password"
            className="form-control"
            id="floatingPassword2"
            placeholder="Confirm password"
            onChange={(event) => { setPassword2(event.target.value) }}
            value={password2}
          />
          <label htmlFor="floatingPassword">Confirm password</label>
        </div>

        <div className="d-grid gap-2 text-center">
          <button
            className="btn btn-primary w-100 py-2"
            type="button"
            onClick={onRegisterButtonClick}
          >
            Register
          </button>
        </div>
      </form>
      <Version />
    </div>
  )
}
