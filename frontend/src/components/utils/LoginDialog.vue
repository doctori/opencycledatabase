<template>
    <div>
        <h2> Login with Strava </h2>
        <a :href="this.getStravaConnectURL()">Connect with Strava</a>
    <div>
        Your token is {{ accessToken }}
    </div>
    </div>
</template>

<script>
import BackendApiClient from '../../services/BackendApiClient'
export default {
    name: 'LoginDialog',
    components: {},
    data(){
        return {
            stravaAPI : {
                clientID : 70293,
                redirectURI : "http://ocd.io:8080/profile"

            },
            accessToken: ""
        }
    },
    mounted() {
        if ( this.$route.query.code){
            console.log("this is the login !")
            this.getToken(this.$route.query.code,this.$route.query.scope,this.$route.query.state)
        }
    },
    methods: {
        getToken(code,scope,state) {
            this.$emit('gerard')
            this.$emit('logged',BackendApiClient.get("/exchange-token",{
                params: {
                    code: code,
                    scope: scope,
                    state: state,
                }
            }))
        },
        getStravaConnectURL(){
            return "http://www.strava.com/oauth/authorize?client_id="+this.stravaAPI.clientID+"&response_type=code&redirect_uri="+this.stravaAPI.redirectURI+"&approval_prompt=force&scope=read"
        }
    }
}
</script>

