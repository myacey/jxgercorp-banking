<template>
  <div class="transfer-block">
    <button class="new-transfer-btn btn h3" @click="$emit('new-transfer')">new transfer</button>
    <h2 class="balance-amount">{{ formattedBalance }}</h2>
    <p class="base-text account-name">{{ username }}</p>
  </div>
</template>

<script>
import { getUserBalance, createTransfer } from '@/api/api';
import { getUsernameFromToken } from '@/utils/auth';
import { snake } from 'case';

export default {
  name: 'NewTransferBlock',

  props: {
    account: {
      type: Object,
      required: true
    }
  },

  data() {
    return {
      username: '',
    }
  },

  computed: {
    currencySymbol() {
      const map = {
        USD: '$',
        EUR: '€',
        RUB: '₽'
      }

      return map[this.account.currency] || this.account.currency;
    },

    formattedBalance() {
      return `${this.account.balance} ${this.currencySymbol}`
    }
  },

  mounted() {
    this.username = getUsernameFromToken()
  },
};
</script>


<style scoped>
.transfer-block {
  align-self: stretch;
  display: flex;
  flex-direction: column;
  align-items: stretch;
  text-align: center;
  justify-content: center;
  width: 130px;
  margin: auto 0;
  gap: 1px;
}

.new-transfer-btn {
  border-radius: 15px;
  display: flex;
  flex: 1;
  align-self: stretch;
  flex-direction: column;
  align-items: stretch;
  justify-content: center;
  background-color: var(--color-accent);
  padding: 4px 13px;

  color: white;
}


.balance-amount {
  color: white;
  font-size: 24px;
  font-weight: 700;
}

.account-name {
  color: var(--color-accent);
}
</style>
