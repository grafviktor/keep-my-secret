/* eslint-disable */
import get from 'lodash/get'
import isNil from 'lodash/isNil'
import { useMemo, useState, useEffect } from 'react'
import Navigation from '../Navigation'
import Error from '../Error'
import Login from '../Login'
import Register from '../Register'
import Home from '../Home'
import PaymentCard from '../Home/PaymentCard'
import File from '../Home/File'
import ApplicationContext from '../../context'
import { isTokenExpired } from '../../utils'
import createAPI from '../../api'

import './style.css'

const viewToComponentMap = {
  login: Login,
  register: Register,
  home: Home,
  card: PaymentCard,
  file: File,
}

export default () => {
  const [alertMessage, setAlertMessage] = useState('')
  const [accessToken, setAccessToken] = useState('')
  const isLoggedIn = () => !isTokenExpired(accessToken)
  const [view, navigateTo] = useState('login')
  const [refreshTokenIngervalID, setRefreshTokenIntervalID] = useState(null)
  const [loggedIn, setLoggedIn] = useState(false)
  const [secret, setSecret] = useState(null)
  const [secrets, setSecrets] = useState({})
  const [appBusy, setAppBusy] = useState(false)

  const api = createAPI({setAppBusy, navigateTo})

  const getVersion = () => api.getVersion()

  const createSecret = async (secret) => {
    return api.createSecret(accessToken, secret)
  }

  const updateSecret = (secret, id) => {
    return api.updateSecret(accessToken, id, secret)
  }

  const deleteSecret = (id) => {
    return api.deleteSecret(accessToken, id)
  }

  const getSecretFile = async (id) => {
    const response = await api.getSecretFile(accessToken, id)
    const url = window.URL.createObjectURL(new Blob([response.data]));
    const link = document.createElement('a');
    const contentDispositionHeader = get(response, 'headers.content-disposition', '')
    const match = contentDispositionHeader.match(/.*filename=([^"]+)/)

    if (match) {
      link.href = url;
      link.setAttribute('download', match[1]);
      document.body.appendChild(link);
      link.click();
    }
  }

  const fetchSecrets = () => {
    return api.fetchSecrets(accessToken)
  }

  useEffect(() => {
    if (!loggedIn) {
      navigateTo('login')
    }
  }, [loggedIn])

  useEffect(() => {
    if (isTokenExpired(accessToken)) {
      setLoggedIn(false)
    } else {
      setLoggedIn(true)
    }
  }, [accessToken])

  useEffect(() => {
    if (!isTokenExpired(accessToken)) {
      setRefreshTokenIntervalID(setInterval(async () => {
        try {
          await api.refreshToken()
        } catch (error) {
          console.warn('Warning: unable to get refreshed token')
        }
      }, 1000 * 60))
    } else {
      clearInterval(refreshTokenIngervalID)
    }
  }, [accessToken])

  const contextValue = useMemo(() => ({
    // verbs
    setSecret,
    getVersion,
    isLoggedIn,
    navigateTo,
    setSecrets,
    createSecret,
    updateSecret,
    fetchSecrets,
    deleteSecret,
    getSecretFile,
    setAccessToken,
    setAlertMessage,

    // and nouns
    api,
    secret,
    secrets,
    loggedIn,
  }))

  const ActiveComponent = viewToComponentMap[view] || Login

  return (
    <div className="application">
      <ApplicationContext.Provider value={contextValue}>
        <Navigation loggedIn={loggedIn} />
        <Error message={alertMessage} view={view} />
        <ActiveComponent loggedIn={loggedIn} appBusy={appBusy} />
      </ApplicationContext.Provider>
    </div>
  )
}
