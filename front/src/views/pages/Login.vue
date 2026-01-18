<template>
  <div class="login-page">
    <div class="login-card">
      <h1 class="login-card__brand brand-name">jxger corp</h1>

      <div class="login-form">
        <h3 class="login-form__title">login</h3>

        <!-- ВВОД USERNAME-->
        <div class="login-form__field">
          <div class="login-form__icon icon-mask" :style="{ maskImage: `url(${userSVG})` }" />
          <input class="login-form__input base-text" id="username" v-model="username" placeholder="username"
            autocomplete="username" />
        </div>

        <!-- ВВОД ПАРОЛЯ -->
        <div class="login-form__field">
          <div class="login-form__icon icon-mask" :style="{ maskImage: `url(${keySVG})` }" />
          <input class="login-form__input base-text" id="password" type="password" v-model="password"
            placeholder="password" autocomplete="current-password" />
        </div>

        <!-- ОШИБКА  -->
        <p v-if="loginError" class="login-form__error">
          {{ loginError }}
        </p>


        <div class="login-form__actions">
          <button class="muted-btn h3" @click="goRegister"> sign up </button>
          <button class="primary-btn h3" @click="handleLogin"> next </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { loginUser } from '@/api/api'

export default {
  data() {
    return {
      username: '',
      password: '',
      loginError: '',

      userSVG: require('@/assets/user.svg'),
      keySVG: require('@/assets/key.svg')
    }
  },

  methods: {
    async handleLogin() {
      // clear errors
      this.loginError = ''

      // check username
      if (!this.username) {
        this.loginError = 'invalid username'
        return
      }

      // check password
      if (!this.password) {
        this.loginError = 'password is too short'
        return
      }

      // call API
      try {
        const userData = {
          username: this.username,
          password: this.password
        }

        await loginUser(userData)

        // clear form
        this.username = ''
        this.password = ''
        this.$router.push('/main') // автоматический переход на /main
      } catch (err) {
        console.error('Login failed:', err)
        this.loginError = err.data?.message
      }
    },
    goRegister() {
      this.$router.push('/register')
    }
  }
}
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: var(--color-bg);
}

.login-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
}

/* ===== form ===== */

.login-form {
  display: flex;
  flex-direction: column;
  gap: 10px;
  min-width: 260px;

  background-color: var(--color-bg-alt);
  border-radius: 10px;
  padding: 8px;
}

.login-form__title {
  text-align: center;
}

.login-form__field {
  display: flex;
  align-items: center;
  gap: 6px;

  background-color: var(--color-bg);
  border-radius: 10px;
  padding: 4px 6px;
}

.login-form__input {
  flex: 1;
  background: none;
  border: none;
  outline: none;
}

/* ===== icon ===== */

.icon-mask {
  width: 17px;
  height: 17px;
  background-color: var(--color-accent);

  mask-size: contain;
  mask-repeat: no-repeat;
  mask-position: center;

  -webkit-mask-size: contain;
  -webkit-mask-repeat: no-repeat;
  -webkit-mask-position: center;
}

/* ===== misc ===== */

.login-form__error {
  text-align: right;
  color: var(--color-accent);
  font-size: 0.9em;
}

.login-form__actions {
  display: flex;
  justify-content: space-between;
  padding-top: 4px;
}
</style>
