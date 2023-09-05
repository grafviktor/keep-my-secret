import {useContext} from 'react'
import Logo from '../Logo'
import ApplicationContext from '../../context'

export default ({loggedIn}) => {
  const {setAccessToken, api} = useContext(ApplicationContext)

  const logout = async () => {
    try {
      await api.logout()
      setAccessToken('')
    } catch (error) {
      console.warn(error)
    }

    // navigateTo('login')
  }

  return (
    <nav className="navbar navbar-expand-lg navbar-light bg-light">
      <div className="container-fluid">
        <button className="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarSupportedContent" aria-controls="navbarSupportedContent" aria-expanded="false" aria-label="Toggle navigation">
          <span className="navbar-toggler-icon" />
        </button>
        <div className="collapse navbar-collapse" id="navbarSupportedContent">
          <ul className="navbar-nav me-auto mb-2 mb-lg-0">
            <li className="nav-item">
              <Logo width="38" height="38" />
            </li>

            <li className="nav-item">
              <a className="nav-link active" aria-current="page" href="/">
                Keep My Secret
              </a>
            </li>
          </ul>
          {loggedIn
            && (
            <div className="d-flex">
              <button className="btn btn-danger" type="button" onClick={logout}>Logout</button>
            </div>
            )}
        </div>
      </div>
    </nav>
  )
}
