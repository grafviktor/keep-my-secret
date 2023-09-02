import cx from 'classnames'
import {useContext, useEffect, useState} from 'react'
import ApplicationContext from '../../context'

import './style.css'

export default ({message = '', view = ''}) => {
  const [isHidden, setHidden] = useState(true)
  const {setAlertMessage} = useContext(ApplicationContext)

  useEffect(() => {
    // There is no point to keep the alert message if we're moving to a different view
    setAlertMessage('')
  }, [view])

  useEffect(() => {
    if (message !== '') {
      setHidden(false)
    } else {
      setHidden(true)
    }
  }, [message])

  const onCloseButtonClick = (event) => {
    event.preventDefault()
    setAlertMessage('')
  }

  return (
    <div
      className={cx(
        'alert alert-danger alert-dismissible kms-alert show fade',
        {hidden: isHidden},
      )}
      role="alert"
    >
      <div>{message}</div>

      <button
        type="button"
        className="btn-close"
        aria-label="Close"
        onClick={onCloseButtonClick}
      />
    </div>
  )
}
