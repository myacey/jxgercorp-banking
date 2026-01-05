<template>
    <div class="account-dropdown" @click="$emit('click')">
        <p class="selector-text">account</p>
        <p>▼</p>
    </div>
</template>

<script>
export default {
    name: 'AccountDropdown'
};
</script>


<style scoped>
.account-dropdown {
    display: flex;
    flex-direction: row;
    gap: 5px;

    cursor: pointer;
}

.account-dropdown:hover {
    opacity: 0.8;
}
</style>