<template>
  <NuxtLayout>
    <v-container class="py-10">
      <v-card rounded="xl" elevation="1" max-width="640">
        <v-card-title class="text-h5 px-6 pt-6">Помилка {{ error.statusCode }}</v-card-title>
        <v-card-text class="px-6 pb-4">
          {{ error.statusCode === 404 ? 'Сторінку не знайдено' : 'Щось пішло не так' }}
        </v-card-text>
        <v-card-actions class="px-6 pb-6">
          <v-btn color="primary" to="/wireguard">На головну</v-btn>
        </v-card-actions>
      </v-card>
    </v-container>
  </NuxtLayout>
</template>

<script setup lang="ts">
defineProps<{
  error: {
    statusCode: number
  }
}>()
</script>
