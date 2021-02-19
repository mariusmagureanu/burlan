import VueRouter from 'vue-router'
import Vue from 'vue'
import App from './App.vue'

import LoginComponent from "./components/login.vue"
import messaging from "./components/messaging.vue"


Vue.config.productionTip = false
Vue.use(VueRouter)

new Vue({
router: new VueRouter({
            routes: [
                {
                    path: '/',
                    redirect: {
                        name: "login"
                    }
                },
                {
                    path: "/login",
                    name: "login",
                    component: LoginComponent
                },
                {
                    path: "/messaging",
                    name: "messaging",
                    component: messaging
                }
            ]
        }),
render: h => h(App),
}).$mount('#app');

