import get from 'lodash/get'
import map from 'lodash/map'
import startCase from 'lodash/startCase'
import React, {useEffect, useState, useContext} from 'react'
import applicationContext from '../../context'

import './style.css'

export default () => {
  const [version, setVersion] = useState({})
  const {getVersion, setAlertMessage} = useContext(applicationContext)

  useEffect(() => {
    (async () => {
      try {
        const {data : response} = await getVersion()
        setVersion(get(response, 'data', {}))
      } catch (error) {
        console.warn(error.message)
        setAlertMessage(`Error: ${error.message}`)
      }
    })()
  }, [])

  return (
    <div className="kms-version">
      <dl>
        {map(version, (value, name) => (
          <React.Fragment key={name}>
            <dt>{startCase(name)}</dt>
            <dd>{value}</dd>
          </React.Fragment>
        ))}
      </dl>
    </div>
  )
}
