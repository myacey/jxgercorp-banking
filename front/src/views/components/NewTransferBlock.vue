<template>
  <div class="transfer-block">
    <div class="action-buttons">
      <div class="new-transfer-btn btn" @click="$emit('new-transfer')">new transfer</div>
    </div>
    <div class="balance-amount">{{ formattedBalance }}</div>
    <div class="account-name">{{ username }}</div>
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
  font: 400 18px EB Garamond, -apple-system, Roboto, Helvetica, sans-serif;
}

.action-buttons {
  align-self: center;
  max-width: 100%;
  width: 130px;
  color: var(--Secondary-Color, #f3f3f3);
}

.new-transfer-btn {
  border-radius: 15px;
  display: flex;
  flex-direction: column;
  align-items: stretch;
  justify-content: center;
  background-color: var(--Primary-Color, #a83232);
  padding: 4px 13px;

  color: var(--Secondary-Color, #f3f3f3);
}


.balance-amount {
  color: #fff;
  font-size: 24px;
  font-weight: 700;
}

.account-name {
  color: #a83232;
}
</style>
