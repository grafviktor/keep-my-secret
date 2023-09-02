import axios from 'axios'
import get from 'lodash/get'
import isNil from 'lodash/isNil'
import reduce from 'lodash/reduce'

const BASE_URL = 'https://localhost:8080'
const httpRequest = axios.create({
  baseURL        : BASE_URL,
  // we handle all responses independently of their HTTP status
  validateStatus : () => true,
})

const checkApiError = (response) => {
  const status = get(response, 'data.status')

  if (response.headers['content-type'] === 'application/json' && status !== 'success') {
    const message = get(response, 'data.message', 'Unknown error')

    throw Error(message)
  } else if (response.status >= 400) {
    throw Error(`status ${response.status}`, response.message)
  }

  return response
}

const getVersion = () => httpRequest.get('/api/v1/version')

const register = (username, password) => httpRequest.post(
  '/api/v1/user/register',
  {username, password},
  {withCredentials: true},
)

const login = (username, password) => httpRequest.post(
  '/api/v1/user/login',
  {username, password},
  {withCredentials: true},
)

const logout = () => httpRequest.post('/api/v1/user/logout')

const refreshToken = () => httpRequest.get(
  '/api/v1/user/token-refresh',
  {withCredentials: true},
)
const createSecretWithFile = async (accessToken, payload) => {
  const {file, ...otherAttributes} = payload
  const formData = new FormData()

  formData.append('data', JSON.stringify({ // Add JSON data
    ...otherAttributes,
    file_name: file.name,
  }))
  formData.append('file', file, file.name) // Add the file

  return httpRequest.post(
    '/api/v1/secrets',
    formData,
    {headers: {
      Authorization  : `Bearer ${accessToken}`,
      'Content-Type' : 'multipart/form-data',
    }},
  )
}

const createSecret = async (accessToken, payload) => {
  if (!isNil(payload.file)) {
    return createSecretWithFile(accessToken, payload)
  }

  return httpRequest.post(
    '/api/v1/secrets',
    payload,
    {headers: {
      Authorization: `Bearer ${accessToken}`,
    }},
  )
}

const updateSecret = (accessToken, id, payload) => httpRequest.put(
  `/api/v1/secrets/${id}`,
  payload,
  {headers: {Authorization: `Bearer ${accessToken}`}},
)

const deleteSecret = (accessToken, id) => httpRequest.delete(
  `/api/v1/secrets/${id}`,
  {headers: {Authorization: `Bearer ${accessToken}`}},
)

const getSecretFile = async (accessToken, id) => httpRequest.get(
  `/api/v1/secrets/file/${id}`,
  {
    // By default responseType is 'json'. If we're requesting binary data(like a file)
    // we receive malformed document, which content is different is different from
    // actual 'Content-Length' header value
    responseType : 'arraybuffer',
    headers      : {Authorization: `Bearer ${accessToken}`},
  },
)

const fetchSecrets = (accessToken) => httpRequest.get(
  '/api/v1/secrets',
  {headers: {Authorization: `Bearer ${accessToken}`}},
)

const api = (context) => reduce({
  // cannot use map[fn1, fn2, ...] because webpack removes function names (fn.name) from
  // prod build don't have time to make a proper setup for WebPack terser plugin
  login,
  logout,
  register,
  getVersion,
  refreshToken,
  createSecret,
  deleteSecret,
  updateSecret,
  fetchSecrets,
  getSecretFile,
}, (acc, fn, name) => ({
  ...acc,
  [name]: async (...args) => {
    context.setAppBusy(true)
    const response = await fn(...args)
    context.setAppBusy(false)

    return checkApiError(response)
  },
}), {})

export default api
