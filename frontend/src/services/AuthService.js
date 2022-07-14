//import { UserManager, WebStorageStateStore } from 'oidc-client'
import {JSO} from 'jso-2'
import BackendApiClient from './BackendApiClient'

let config = {
	providerID: "strava",
	client_id: "70293",
	redirect_uri: "http://ocd.io:8080/profile", // The URL where you is redirected back, and where you perform run the callback() function.
	authorization: "https://www.strava.com/oauth/authorize",
    response_type: "code",
	scopes: { request: ["read"]}
}
let j = new JSO(config)
/**
 * Class to encapsulate all authentication related logic.
 */
class AuthService {
    token = {}
    /**
     * Initate the login process.
     */
    login () {
        console.log("honettement, on est où là ? ")
        j.getToken()
        .catch(error => console.log(error))
    }

    logout () {
        j.wipeTokens()
        .then(() => console.log('User logged out'))
        .catch(error => console.log(error))
    }

    /**
     * Handles the redirect from the OAuth server after a user logged in.
     */
    async handleLoginRedirect (data) {
        if (Object.prototype.hasOwnProperty.call(data,'state')){
            j.callback(data)
            const response = await BackendApiClient.get("/exchange-token", {
                params: {
                    code: data.code,
                    scope: data.scope,
                    state: data.state,
                }
            })
            this.token = response.data
            localStorage.setItem("tokens-strava", JSON.stringify(this.token))  
        }
   
    }

    /**
     * Handles the redirect from the OAuth server after a user logged out.
     */
    handleLogoutRedirect () {
        return j.signoutRedirectCallback()
    }

    /**
     * Checks whether or not a user is currently logged in.
     *
     * Returns a promise which will be resolved to true/false or be rejected with an error.
     */
    isUserLoggedIn () {
        return new Promise((resolve) => {
            this.token = localStorage.getItem('tokens-strava')
            if (this.token === null) {
                resolve(false)
            }
                resolve(true)
        })
    }

    /**
     * Get the profile data for the currently authenticated user.
     *
     * Returns an empty object if no user is logged in.
     */
    getProfile () {
        return new Promise((resolve, reject) => {
        j.getUser()
            .then(user => {
            if (user === null) {
                resolve(null)
            }
            resolve(user.profile)
            })
            .catch(error => reject(error))
        })
    }

    /**
     * Get the access token.
     *
     * Can be used to make requests to the backend.
     */
    getAccessToken () {
        return new Promise((resolve, reject) => {
            console.log('Get access token from user')
            j.getUser()
            .then(user => {
                console.log('Got access token from user')
                resolve(user.access_token)
            })
            .catch(error => reject(error))
        })
    }
}

/**
 * Create and expose an instance of the auth service.
 */
 export const authService = new AuthService()

 /**
  * Default export to register the authentication service in the global Vue instance.
  *
  * This allows us to reference it using "this.$auth" whenever we are inside of a Vue context.
  */
 export default {
   install: function (Vue) {
     Vue.prototype.$auth = authService
   }
 }