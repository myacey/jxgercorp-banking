<template>
  <div class="transaction-item">
    <div :class="['amount', type]">
      {{ type == 'positive' ? '+' : '-' }}{{ amount }}
    </div>
    <div class="transaction-details">
      <p class="user-name">{{ user }}</p>
      <p class="user-date-divider"> | </p>
      <p class="very-small-text transaction-time ">{{ time }}</p>
    </div>
  </div>
</template>

<script>
export default {
  name: 'TransactionItem',
  props: {
    amount: {
      type: Number,
      required: false,
      default: '0.00'
    },
    user: {
      type: String,
      required: false,
      default: 'unknown'
    },
    time: {
      type: String,
      required: false,
      default: new Date().toLocaleTimeString()
    },
    type: {
      type: String,
      required: false,
      validator: (value) => ['positive', 'negative'].includes(value),
      default: 'positive'
    }
  }
}
</script>

<style scoped>
.transaction-item {
  display: flex;
  flex: 1;
  align-items: center;
  justify-content: space-between;
  font-family: EB Garamond, -apple-system, Roboto, Helvetica, sans-serif;
}

.amount {
  font-size: 16px;
  align-self: stretch;
  margin: auto 0;
  margin-right: 8px;
}

.amount.positive {
  color: #008414;
}

.amount.negative {
  color: #a83232;
}

.transaction-details {
  display: flex;
  flex-direction: row;
  flex: 1;
  margin-left: 8px;
  gap: 2px;
  align-items: stretch;
  justify-content: right;
  margin-right: 4px;
}


.user-name {
  font-size: 16px;
  align-self: center;
  min-width: 0px;
  white-space: nowrap;
}

.user-date-divider {
  align-self: center;
  white-space: nowrap;

  color: var(--color-text-muted);
}

.transaction-time {
  color: var(--color-text-muted);
}
</style>
