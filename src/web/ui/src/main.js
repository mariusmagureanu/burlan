import VueRouter from 'vue-router'
import Vue from 'vue'
import App from './App.vue'

import LoginComponent from "./components/login.vue"
import SecureComponent from "./components/secure.vue"
import UsersList from "./components/UsersList.vue"

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
                    path: "/secure",
                    name: "secure",
                    component: SecureComponent
                },
                {
                    path: "/users",
                    name: "users",
                    component: UsersList
                }
            ]
        }),
render: h => h(App),
}).$mount('#app');

