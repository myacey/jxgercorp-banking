<template>
    <div class="modal">
        <h2 class="header">delete account?</h2>
        <p class="warn-text">this action cannot be undone!</p>

        <div class="account-id-data">
            <span class="acctound-id-text base-text">
                account id:
            </span>
            <span class="account-id">
                {{ account.id }}
            </span>
        </div>

        <div class="btns">
            <h3 class="cancel-btn btn" @click="$emit('cancel')"> cancel </h3>
            <h3 class="delete-btn btn" @click="$emit('confirm', account)"> delete </h3>
        </div>
    </div>
</template>

<script>
export default {
    name: "DeleteAccountModal",
    props: {
        account: {
            type: Object,
            required: true,
        }
    },
}
</script>



<style scoped>
.modal {
    display: flex;
    flex: 1;
    max-width: fit-content;
    flex-direction: column;
    align-items: stretch;

    background-color: var(--color-bg-alt);
    padding: 4px;
    gap: 10px;
    border-radius: 10px;
    justify-content: center;
}

.header {
    flex: 1;
    align-self: stretch;
    text-align: center;
}

.warn-text {
    flex: 1;
    align-self: stretch;

    text-align: end;
    color: var(--color-accent)
}

.account-id-data {
    display: flex;
    flex: 1;
    align-self: stretch;
    flex-direction: row;
    justify-content: space-between;
    gap: 10px;
    padding: 4px;
}

.account-id-text {
    text-align: left;
}

.account-id {
    flex: 1;
    align-self: stretch;
    text-align: right;
    color: var(--color-bg)
}

.btns {
    display: flex;
    flex: 1;
    flex-direction: row;
    align-self: stretch;
    justify-content: space-between;
    padding: 5px;
}

.cancel-btn {
    background-color: var(--color-accent);
    border-radius: 10px;
    text-align: center;
    align-items: center;
    color: white;
    padding: 1px 7px;
}

.delete-btn {
    border: 1px solid var(--color-text-muted);
    border-radius: 10px;
    text-align: center;
    align-items: center;
    color: var(--color-text-muted);
    padding: 1px 7px;
}
</style>
