<template>
  <nav class="navigation">
    <div class="navigation__brand">
      <h1 class="navigation__brand-name brand-name">jxger corp</h1>
    </div>

    <div class="navigation__menu">
      <a href="#" class="nav-link">home</a>
      <a href="#" class="nav-link">features</a>
      <a href="#" class="nav-link">about us</a>
      <a href="#" class="nav-link">contact</a>

      <span class="navigation__separator">|</span>
      <div class="navigation__user">
        <div class="avatar-container__icon icon-mask" :style="{ maskImage: `url(${userSVG})` }" />
        <h3 class="navigation__username">{{ username ? username : 'username' }}</h3>
      </div>
    </div>
  </nav>
</template>

<script>
import { getUsernameFromToken } from '@/utils/auth';

export default {
  name: 'AppNavigation',
  data() {
    return {
      username: '',
      userSVG: require('@/assets/user.svg'),
    }
  },
  mounted() {
    this.username = getUsernameFromToken();
  }
}
</script>

<style scoped>
.navigation {
  display: flex;
  align-items: center;
  gap: 40px 100px;
  justify-content: space-between;
  flex-wrap: wrap;
}

.navigation__brand {
  align-self: stretch;
  display: flex;
  min-height: 35px;
  flex-direction: column;
  overflow: hidden;
  align-items: center;
  color: white;
  text-align: center;
  justify-content: center;
  margin: auto 0;
  padding: 3px 10px;
  font: 400 40px Runic AltNo, -apple-system, Roboto, Helvetica, sans-serif;
}

.brand-name {
  color: white;
  margin: 0;
  font-size: inherit;
  font-weight: inherit;
  font-family: inherit;
}

.navigation__menu {
  align-self: stretch;
  display: flex;
  min-width: 240px;
  align-items: center;
  gap: 14px;
  justify-content: center;
  margin: auto 0;
  padding: 0 5px;
}

.nav-link {
  color: #f3f3f3;
  text-decoration: none;
  text-align: center;
  align-self: stretch;
  margin: auto 0;
  font: 400 18px EB Garamond, -apple-system, Roboto, Helvetica, sans-serif;
}

.nav-link:hover {
  opacity: 0.8;
}

.navigation__separator {
  color: #f3f3f3;
  text-align: center;
  align-self: stretch;
  margin: auto 0;
  font: 400 18px EB Garamond, -apple-system, Roboto, Helvetica, sans-serif;
}

.navigation__user {
  display: inline-flex;
  align-items: center;
  gap: 15px;

  border: 1px solid var(--color-text-muted);
  border-radius: 15px;
  padding: 2px 6px;

}

.navigation__username {
  color: var(--color-accent);
}

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
</style>
