<template>
    <div>
        <div v-if="profile == null">
            <button type="button" class="login-button button"  v-if="!isUserLoggedIn" @click="onLogin()">Login with Strava</button>  
            <button v-if="isUserLoggedIn" @click="onLogout()">Logout</button>  
        </div>
        <div v-if="profile !== null">
            <h2>You are : {{ user.login }} </h2>
            <div>
                Your token is {{ user.token }} <br/>
                Country : {{ user.country }}
            </div>
        </div>
    </div>
</template>

<script>
//import BackendApiClient from '../../services/BackendApiClient'
import { registerUserLoggedInEventListener, registerUserLoggedOutEventListener } from '@/services/EventBus'
export default {
    name: 'ProfileInfo',
    components: {},
    data () {
        return {
            isUserLoggedIn: false,
            profile: null
        }
    },
    mounted () {
        console.log(this.$route.query)
        this.$auth.handleLoginRedirect(this.$route.query)


        this.$auth.isUserLoggedIn()
            .then(isLoggedIn =>{
                this.isUserLoggedIn = isLoggedIn
            })
            .catch(error => {
                console.log(error)
                this.isUserLoggedIn = false
            })
        registerUserLoggedInEventListener(() => { this.isUserLoggedIn = true })
        registerUserLoggedOutEventListener(() => { this.isUserLoggedIn = false })
        if (this.isUserLoggedIn){
            this.$auth.getProfile()
                .then(profile => {
                    this.profile = profile
                })
                .catch(error => {
                    console.log(error)
                    this.profile = {}
                })
        }
    },
    methods: {
        onLogin(){
            console.log(this.$auth)
            this.$auth.login()
        },
        onLogout(){
            this.$auth.logout()
        },
    }
}
</script>

