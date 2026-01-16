<template>
  <div class="transaction-history">
    <h2 class="history-title">history</h2>
    <div class="history-list">
      <TransactionDay v-for="(dayTransactions, date) in groupedTransactions" :key="date" :date="date"
        :transactions="dayTransactions" />
    </div>
  </div>
</template>

<script>
import TransactionDay from './TransactionDay.vue'
import { convertKeysToCamel } from '@/utils/snake2camel';
import { snake } from 'case';
import { fetchTransfers } from '@/api/api';

export default {
  name: 'TransactionHistory',

  props: {
    accountID: {
      type: String,
      required: true
    },
    transactions: {
      type: Array,
      required: false,
    }
  },

  components: {
    TransactionDay
  },

  computed: {
    groupedTransactions() {
      const groups = {}

      this.transactions.forEach(trx => {
        const date = new Date(trx.createdAt).toLocaleDateString('ru-RU')

        if (!groups[date]) {
          groups[date] = []
        }

        groups[date].push({
          amount: trx.amount,
          username: this.resolveUser(trx),
          time: new Date(trx.createdAt).toLocaleString([], {
            hour: '2-digit',
            minute: '2-digit'
          }),
          type: trx.fromAccountId == this.accountID ? 'negative' : 'positive'
        })
      })

      return groups
    }
  },

  methods: {
    resolveUser(trx) {
      return trx.fromAccountId == this.accountID ? trx.toAccountUsername : trx.fromAccountUsername
    }
  }
}
</script>

<style scoped>
.transaction-history {
  border-radius: 10px;
  align-self: stretch;
  flex-direction: column;
  align-items: stretch;
  justify-content: center;
  margin: auto 0;
  border: 1px solid rgba(40, 40, 40, 1);

  min-height: 340px;
  min-width: 220px;
}

.history-title {
  color: #fff;
  flex: 1;
  align-self: stretch;
  text-align: center;
  font: 400 18px EB Garamond, -apple-system, Roboto, Helvetica, sans-serif;
  margin: 0;
}

.history-list {
  margin-top: 5px;
  flex: 1;
}
</style>
