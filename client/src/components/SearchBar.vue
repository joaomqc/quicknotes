<template>
  <div>
    <div class="rounded-md border px-5 py-4">
      <input
        type="text"
        ref="input"
        v-model="searchTerm"
        v-focus
        class="w-full bg-transparent focus:outline-none"
        placeholder="Search..."
        @keydown="keydownHandler"
      />
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const emit = defineEmits(['search'])

const input = ref()
const searchTerm = ref()
const router = useRouter()

function keydownHandler(event) {
  if (event.key === 'Enter') {
    search()
  }
}
function search() {
  router.push({
    name: 'search',
    query: { ['q']: searchTerm.value },
  })
  emit('search')
}

</script>
