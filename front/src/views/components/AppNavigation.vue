<template>
  <nav class="navigation">
    <div class="logo">
      <div class="logo-text">
        <h1 class="brand-name">jxger corp</h1>
      </div>
    </div>

    <div class="menu">
      <a href="#" class="nav-link">home</a>
      <a href="#" class="nav-link">features</a>
      <a href="#" class="nav-link">about us</a>
      <a href="#" class="nav-link">contact</a>

      <span class="seoarator">|</span>
      <div class="avatar-container">
        <div class="avatar-container">
          <img
            src="https://api.builder.io/api/v1/image/assets/TEMP/4724a1e08316ed1696b11e825015f45acbfa5730?placeholderIfAbsent=true&apiKey=2273adaceecb4f3fb0db700e4dc8777b"
            alt="User avatar" class="avatar-image" />
          <div class="avatar-border"></div>
        </div>
      </div>
    </div>
  </nav>
</template>

<script>
export default {
  name: 'AppNavigation'
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

@media (max-width: 991px) {
  .navigation {
    max-width: 100%;
  }
}

.logo {
  align-self: stretch;
  display: flex;
  min-height: 35px;
  flex-direction: column;
  overflow: hidden;
  align-items: center;
  color: rgba(243, 243, 243, 1);
  text-align: center;
  justify-content: center;
  margin: auto 0;
  padding: 3px 10px;
  font: 400 40px Runic AltNo, -apple-system, Roboto, Helvetica, sans-serif;
}

.logo-text {
  max-width: 100%;
  width: 125px;
  padding: 0 4px;
}

.brand-name {
  color: white;
  margin: 0;
  font-size: inherit;
  font-weight: inherit;
  font-family: inherit;
}

.menu {
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

.separator {
  color: #f3f3f3;
  text-align: center;
  align-self: stretch;
  margin: auto 0;
  font: 400 18px EB Garamond, -apple-system, Roboto, Helvetica, sans-serif;
}

.user-avatar {
  border-radius: 25px;
  align-self: stretch;
  display: flex;
  min-height: 29px;
  align-items: center;
  justify-content: space-between;
  width: 125px;
  margin: auto 0;
  padding: 1px 5px;
}

.avatar-container {
  display: flex;
  flex-direction: column;
  border-radius: 0px 0px 0px 0px;
  align-self: stretch;
  position: relative;
  aspect-ratio: 4.259;
  width: 115px;
  margin: auto 0;
}

.avatar-image {
  position: absolute;
  inset: 0;
  height: 100%;
  width: 100%;
  object-fit: cover;
  object-position: center;
}

.avatar-border {
  position: relative;
  border-radius: 17px;
  z-index: 10;
  display: flex;
  flex-shrink: 0;
  height: 29px;
  border: 1px solid rgba(243, 243, 243, 1);
}
</style>
