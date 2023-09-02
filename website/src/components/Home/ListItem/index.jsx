import {useContext} from 'react'
import ApplicationContext from '../../../context'

import './style.css'

export default ({secret}) => {
  const {
    secrets,
    setSecret,
    navigateTo,
    setSecrets,
    deleteSecret,
    setAlertMessage,
  } = useContext(ApplicationContext)

  const onOpenButtonClick = () => {
    setSecret(secret)
    navigateTo(secret.type)
  }

  const onDeleteButtonClick = async () => {
    try {
      const {data : response} = await deleteSecret(secret.id)
      const {[response.data] : _, ...other} = secrets

      setSecrets(other)
    } catch (error) {
      console.warn(error)
      setAlertMessage(`Error: ${error.message}`)
    }
  }

  const getType = (type) => ({
    card : 'Bank Card',
    file : 'Secret File',
  }[type] || type)

  return (
    <div className="kms-secret-list__item">
      <div className="card kms-secret-card">
        <div className="card-body">
          <h5 className="card-title">{secret.title}</h5>
          <p className="card-text">{getType(secret.type)}</p>
          <div className="kms-secret-card__control-pane">
            <button
              className="btn btn-sm btn-primary"
              type="button"
              onClick={onOpenButtonClick}
            >
              Open
            </button>

            <button
              className="btn btn-sm btn-danger"
              type="button"
              onClick={onDeleteButtonClick}
            >
              Delete
            </button>
          </div>
        </div>
      </div>
    </div>
  )
}
