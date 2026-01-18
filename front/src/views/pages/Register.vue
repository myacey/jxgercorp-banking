<template>
  <div class="register-page">
    <div class="register-card">
      <h1 class="register-card__brand brand-name">jxger corp</h1>

      <form v-if="!isSuccess" class="register-form" @submit.prevent="handleRegister" novalidate>
        <h3 class="register-form__title">register</h3>

        <!-- username -->
        <div class="register-form__field">
          <div class="register-form__icon icon-mask" :style="{ maskImage: `url(${userSVG})` }" />
          <input class="register-form__input base-text" id="username" v-model="username" placeholder="username"
            autocomplete="username" required />
        </div>

        <!-- email -->
        <div class="register-form__field">
          <div class="register-form__icon icon-mask" :style="{ maskImage: `url(${emailSVG})` }" />
          <input class="register-form__input base-text" id="email" v-model="email" placeholder="email"
            autocomplete="email" required />
        </div>

        <!-- password -->
        <div class="register-form__field">
          <div class="register-form__icon icon-mask" :style="{ maskImage: `url(${keySVG})` }" />
          <input class="register-form__input base-text" id="password" v-model="password" placeholder="password"
            type="password" required autocomplete="new-password" />
        </div>

        <!-- confirm password -->
        <div class="register-form__field">
          <div class="register-form__icon icon-mask" :style="{ maskImage: `url(${keySVG})` }" />
          <input class="register-form__input base-text" id="current-password" v-model="confirmPassword"
            placeholder="confirm password" type="password" required autocomplete="new-password" />
        </div>

        <p v-if="registerError" class="register-form__error">{{ registerError }}</p>

        <div class="register-form__actions">
          <button class="muted-btn h3" @click="goLogin">login</button>
          <button class="primary-btn h3" @click="handleRegister" type="submit">sign up</button>
        </div>
      </form>

      <div v-else class="register-success">
        <h3 class="register-success_title">check you email</h3>

        <p class="register-success_text base-text">
          we've sent confirm link to {{ email }}
        </p>

        <p class="register-success__hint small-text">
          Please follow the link to activate your account.
        </p>


        <button class="primary-btn h3" @click="$router.push('/login')">go to login</button>
      </div>
    </div>
  </div>
</template>

<script>
import { registerUser } from '@/api/api'

const USERNAME_REGEX = /^[a-zA-Z0-9_]*$/
const EMAIL_REGEX = /^[^\s@]+@[^\s@]+\.[^\s@]+$/

export default {
  data() {
    return {
      username: '',
      email: '',
      password: '',
      confirmPassword: '',
      registerError: '',

      isSuccess: false,

      userSVG: require('@/assets/user.svg'),
      emailSVG: require('@/assets/email.svg'),
      keySVG: require('@/assets/key.svg'),
    }
  },

  methods: {
    async handleRegister() {
      // clear errors
      this.registerError = ''

      if (!USERNAME_REGEX.test(this.username) || !this.username) {
        console.log('username can contain only letters, numbers and underscore')
        this.registerError = 'username can contain only letters, numbers and underscore'
        return
      }

      // validate email
      if (!EMAIL_REGEX.test(this.email)) {
        this.registerError = 'Invalid email format'
        return
      }

      // check password
      if (this.password !== this.confirmPassword) {
        this.registerError = 'Passwords do not match'
        return
      }

      if (this.password.length < 8) {
        this.registerError = 'Password too short'
        return
      }

      // call API
      try {
        const userData = {
          email: this.email,
          username: this.username,
          password: this.password
        }

        await registerUser(userData)

        this.isSuccess = true
      } catch (err) {
        console.log('Registration failed:', err)
        this.registerError =
          err.data?.message || 'Registration failed. Please try again'
        console.log('Extracted err message:', this.error)
      }
    },

    goLogin() {
      this.$router.push('/login')
    }
  }
}
</script>


<style scoped>
.register-page {
  min-height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: var(--color-bg);
}

.register-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
}

/* ===== form ===== */

.register-form {
  display: flex;
  flex-direction: column;
  gap: 10px;
  min-width: 260px;
  max-width: 360px;

  background-color: var(--color-bg-alt);
  border-radius: 10px;
  padding: 8px;
}

.register-form__title {
  text-align: center;
}

.register-form__field {
  display: flex;
  align-items: center;
  gap: 6px;

  background-color: var(--color-bg);
  border-radius: 10px;
  padding: 4px 6px;
}

.register-form__input {
  flex: 1;
  background: none;
  border: none;
  outline: none;
}

.register-success {
  display: flex;
  flex-direction: column;
  gap: 10px;
  min-width: 260px;
  max-width: 360px;

  background-color: var(--color-bg-alt);
  border-radius: 10px;
  padding: 8px;
}


.register-success_title {
  text-align: center;
}

.register-success__hint {
  opacity: 0.8;
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

.register-form__error {
  text-align: right;
  color: var(--color-accent);

  word-break: break-word;
  overflow-wrap: anywhere;
}

.register-form__actions {
  display: flex;
  justify-content: space-between;
  padding-top: 4px;
}
</style>
