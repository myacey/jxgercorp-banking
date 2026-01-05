import Vue from "vue";
import App from "./App.vue";
import "./assets/styles/buttons.css";
import "./assets/styles/fields.css";
import "./assets/styles/titles.css";
import "./assets/styles/globals.css";
import "./assets/styles/colors.css";
import router from "./router";

Vue.config.productionTip = false;

new Vue({
  router,
  render: (h) => h(App),
}).$mount("#app");
