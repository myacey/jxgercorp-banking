<template>
  <div class="day-list">
    <div class="date-divider">
      <div class="divider-line">
        <div class="line"></div>
      </div>
      <time class="date-text">{{ date }}</time>
      <img
        src="https://api.builder.io/api/v1/image/assets/TEMP/e0a81dcb89dc335d1d329de8287f1a186ae34aa6?placeholderIfAbsent=true&apiKey=2273adaceecb4f3fb0db700e4dc8777b"
        alt="Date indicator" class="date-icon" />
    </div>
    <TransactionItem v-for="(transaction, index) in transactions" :key="index" :amount="transaction.amount"
      :user="transaction.user" :time="transaction.time" :type="transaction.type" />
    <div v-if="showSeeMore" class="see-more-container">
      <button class="see-more-btn">see more</button>
      <img
        src="https://api.builder.io/api/v1/image/assets/TEMP/c6095af68f6049f653281e45e00ddc29076c5222?placeholderIfAbsent=true&apiKey=2273adaceecb4f3fb0db700e4dc8777b"
        alt="Expand" class="expand-icon" />
    </div>
  </div>
</template>

<script>
import TransactionItem from './TransactionItem.vue'

export default {
  name: 'TransactionDay',
  components: {
    TransactionItem
  },
  props: {
    date: {
      type: String,
      required: true
    },
    transactions: {
      type: Array,
      required: true
    },
    showSeeMore: {
      type: Boolean,
      default: false
    }
  }
}
</script>

<style scoped>
.day-list {
  display: flex;
  margin-top: 10px;
  flex: 1;
  flex-direction: column;
  overflow: hidden;
  align-items: stretch;
  justify-content: start;
  padding: 0 10px 10px;
}

.date-divider {
  display: flex;
  flex: 1;
  align-items: center;
  gap: 1px;
  justify-content: center;
  padding: 2px;
}

.divider-line {
  align-self: stretch;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  align-items: stretch;
  justify-content: center;
  flex: 1;
  flex-shrink: 1;
  flex-basis: 0%;
  margin: auto 0;
}

.line {
  border-radius: 5px;
  background-color: rgba(217, 217, 217, 1);
  display: flex;
  min-height: 1px;
  width: 100%;
}

.date-text {
  color: #fff;
  text-align: center;
  align-self: stretch;
  margin: auto 0;
  font: 400 12px EB Garamond, -apple-system, Roboto, Helvetica, sans-serif;
}

.date-icon {
  aspect-ratio: 1;
  object-fit: contain;
  object-position: center;
  width: 7px;
  align-self: stretch;
  flex-shrink: 0;
  margin: auto 0;
}

.see-more-container {
  align-self: center;
  display: flex;
  align-items: center;
  gap: 1px;
  color: #fff;
  text-align: center;
  justify-content: center;
  font: 400 16px EB Garamond, -apple-system, Roboto, Helvetica, sans-serif;
}

.see-more-btn {
  align-self: stretch;
  margin: auto 0;
  background: none;
  border: none;
  color: inherit;
  font: inherit;
  cursor: pointer;
  padding: 0;
}

.see-more-btn:hover {
  opacity: 0.8;
}

.expand-icon {
  aspect-ratio: 1;
  object-fit: contain;
  object-position: center;
  width: 19px;
  align-self: stretch;
  flex-shrink: 0;
  margin: auto 0;
}
</style>
