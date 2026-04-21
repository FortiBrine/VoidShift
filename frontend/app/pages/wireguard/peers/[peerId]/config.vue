<template>
  <div class="mb-6 d-flex align-center justify-space-between">
    <v-btn variant="text" prepend-icon="mdi-arrow-left" :to="backPath">Назад</v-btn>
    <v-btn color="primary" variant="outlined" prepend-icon="mdi-download" @click="downloadConfig">
      Завантажити конфіг
    </v-btn>
  </div>

  <v-card rounded="xl" elevation="1">
    <v-card-title class="text-h5 pt-6 px-6">Конфіг peer-а #{{ peerId }}</v-card-title>
    <v-card-text class="px-6 pb-6">
      <v-textarea
        v-if="!loading"
        :model-value="peerConfig"
        variant="outlined"
        auto-grow
        rows="16"
        readonly
      />

      <div v-else class="py-10 text-center">
        <v-progress-circular indeterminate color="primary" />
      </div>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
const route = useRoute()
const router = useRouter()
const wireguardApi = useWireguardApi()
const api = useApiClient()
const notification = useNotification()

const loading = ref(true)
const peerConfig = ref('')

const peerId = computed(() => Number(route.params.peerId))
const networkIdQuery = computed(() => Number(route.query.networkId))

const backPath = computed(() => {
  if (Number.isInteger(networkIdQuery.value) && networkIdQuery.value > 0) {
    return `/wireguard/networks/${networkIdQuery.value}`
  }

  return '/wireguard'
})

if (!Number.isInteger(peerId.value) || peerId.value <= 0) {
  notification.showError('Некоректний ідентифікатор peer-а')
  await router.push('/wireguard')
}

const loadConfig = async () => {
  loading.value = true

  try {
    peerConfig.value = await wireguardApi.getPeerConfig(peerId.value)
  } catch (error) {
    notification.showError(api.getErrorMessage(error, 'Не вдалося отримати конфіг'))
  } finally {
    loading.value = false
  }
}

const downloadConfig = () => {
  window.open(`/api/vpn/wireguard/peers/${peerId.value}/config/download`, '_blank', 'noopener')
}

await loadConfig()
</script>
