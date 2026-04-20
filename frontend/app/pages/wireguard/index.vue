<template>
  <div class="d-flex align-center justify-space-between mb-6">
    <div>
      <h1 class="text-h4 font-weight-medium">WireGuard</h1>
      <p class="text-medium-emphasis">Керування мережами та peer-конфігами</p>
    </div>

    <v-btn
      color="primary"
      prepend-icon="mdi-plus"
      to="/wireguard/create-network"
    >
      Створити мережу
    </v-btn>
  </div>

  <v-card rounded="xl" elevation="1">
    <v-list v-if="networks.length > 0" lines="two">
      <v-list-item
        v-for="network in networks"
        :key="network.id"
        :title="network.name"
        :subtitle="`${network.address} · Порт ${network.listen_port}`"
        :to="`/wireguard/networks/${network.id}`"
        prepend-icon="mdi-lan"
        append-icon="mdi-chevron-right"
      />
    </v-list>

    <v-card-text v-else-if="!loading" class="py-10 text-center text-medium-emphasis">
      Мереж поки немає. Створи першу WireGuard-мережу.
    </v-card-text>

    <v-card-text v-if="loading" class="py-10 text-center">
      <v-progress-circular indeterminate color="primary" />
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import type { NetworkSummary } from '~/composables/useWireguardApi'

const wireguardApi = useWireguardApi()
const api = useApiClient()
const notification = useNotification()

const loading = ref(true)
const networks = ref<NetworkSummary[]>([])

const loadNetworks = async () => {
  loading.value = true

  try {
    networks.value = await wireguardApi.listNetworks()
  } catch (error) {
    notification.showError(api.getErrorMessage(error, 'Не вдалося завантажити мережі'))
  } finally {
    loading.value = false
  }
}

await loadNetworks()
</script>
