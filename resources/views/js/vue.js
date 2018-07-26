
import Vue from 'vue'
import App from './../Components/App.vue'
import VueMaterial from 'vue-material'
import 'vue-material/dist/vue-material.min.css'

Vue.use(VueMaterial);

new Vue({
    el: '#app',
    render: h => h(App)
});
