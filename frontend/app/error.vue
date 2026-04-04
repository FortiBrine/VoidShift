<template>
  <NuxtLayout>
    <h1 class="text-2xl">{{ error.statusCode }}</h1>
    <p class="text-2xl" v-if="error.statusCode === 404">
      Сторінку не знайдено
    </p>
  </NuxtLayout>
</template>

<script>
export default {
  props: ['error']
};
</script>