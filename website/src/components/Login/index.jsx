import every from 'lodash/every'
import {useContext, useState, useEffect} from 'react'
import Logo from '../Logo'
import Version from '../Version'
import ApplicationContext from '../../context'
import './style.css'

export default ({loggedIn}) => {
  const [username, setUsername] = useState('user@localhost')
  const [password, setPassword] = useState('12345')
  const {setAlertMessage, setAccessToken, navigateTo, api} = useContext(ApplicationContext)

  useEffect(() => {
    if (loggedIn) {
      navigateTo('home')
    }
  }, [loggedIn])

  const onLoginButtonClick = async (ev) => {
    ev.preventDefault()

    if (!every([username, password])) {
      setAlertMessage('Username or password cannot be empty')

      return
    }

    try {
      const {data: response} = await api.login(username, password)

      setAccessToken(response.data)
    } catch (error) {
      console.warn(error)
      setAlertMessage(`Error: ${error.message}`)
    }
  }

  const onRegisterClick = (ev) => {
    ev.preventDefault()

    navigateTo('register')
  }

  return (
    <div className="form-signin w-100 m-auto">
      <Logo />
      <form>
        <h1 className="h3 mb-3 fw-normal">
          Please sign in or
          {' '}
          <a href="register" onClick={onRegisterClick}>register</a>
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

        <div className="d-grid gap-2 text-center">
          <button
            className="btn btn-primary w-100 py-2"
            type="button"
            onClick={onLoginButtonClick}
          >
            Sign in
          </button>
        </div>
      </form>
      <Version />
    </div>
  )
}
