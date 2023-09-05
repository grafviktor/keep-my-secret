import cx from 'classnames'
import get from 'lodash/get'
import map from 'lodash/map'
import {useContext, useEffect} from 'react'
import ApplicationContext from '../../context'
import SecretItem from './ListItem'
import './style.css'

export default ({appBusy}) => {
  const {
    navigateTo,
    fetchSecrets,
    setAlertMessage,
    setSecret,
    setSecrets,
    secrets,
  } = useContext(ApplicationContext)

  useEffect(() => {
    (async () => {
      try {
        const {data: response} = await fetchSecrets()

        setSecrets(response.data)
      } catch (error) {
        console.warn(error.message)
        setAlertMessage(`Error: ${error.message}`)
      }
    })()
  }, [])

  const onNewItemButtonClick = (event) => {
    setSecret(null)

    const view = get(event, 'target.id', 'home')

    navigateTo(view)
  }

  return (
    <div className="kms-home">
      <div className="kms-secret-list">
        {map(secrets, (secret, id) => <SecretItem key={id} secret={{id, ...secret}} />)}
        <div className={cx(
          'spinner-grow text-primary kms-busy-indicator',
          {'kms-busy-indicator__shown': appBusy},
        )}
        >
          <span className="visually-hidden">Loading...</span>
        </div>
      </div>
      <div className="btn-group dropup kms-button-add">
        <button type="button" className="btn btn-primary dropdown-toggle" data-bs-toggle="dropdown" aria-expanded="false">
          New Secret
        </button>
        <ul className="dropdown-menu">
          <li>
            <button
              type="button"
              className="dropdown-item"
              id="card"
              onClick={onNewItemButtonClick}
            >
              Payment card
            </button>
          </li>

          <li>
            <button
              type="button"
              className="dropdown-item"
              id="pass"
              onClick={onNewItemButtonClick}
            >
              Password
            </button>
          </li>

          <li>
            <button
              type="button"
              className="dropdown-item"
              id="note"
              onClick={onNewItemButtonClick}
            >
              Note
            </button>
          </li>

          <li>
            <button
              type="button"
              className="dropdown-item"
              id="file"
              onClick={onNewItemButtonClick}
            >
              File
            </button>
          </li>
        </ul>
      </div>
    </div>
  )
}
