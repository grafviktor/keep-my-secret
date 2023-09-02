/* eslint-disable import/prefer-default-export */

import jwtDecode from 'jwt-decode'

export const isTokenExpired = (token) => {
  try {
    const decodedToken = jwtDecode(token)
    const currentTime = Date.now() / 1000 // Convert to seconds

    // Compare the expiration time with the current time
    return decodedToken.exp < currentTime
  } catch (error) {
    // Handle decoding errors
    return true // Assume token is expired if there's an error
  }
}
