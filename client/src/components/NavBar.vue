<template>
    <nav class="ml-6 mb-2 flex align-top">
        <RouterLink class="hover:bg-inherit" :to="{name: 'home'}">
            <div class="font-medium italic text-3xl font-serif">
                quicknotes
            </div>
        </RouterLink>
    </nav>
</template>

<script setup>
import {RouterLink} from "vue-router";

</script>
