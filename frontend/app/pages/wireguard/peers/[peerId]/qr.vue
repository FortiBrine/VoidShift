<template>
  <div class="mb-6">
    <v-btn variant="text" prepend-icon="mdi-arrow-left" :to="backPath">Назад</v-btn>
  </div>

  <v-card rounded="xl" elevation="1" max-width="760">
    <v-card-title class="text-h5 pt-6 px-6">QR peer-а #{{ peerId }}</v-card-title>
    <v-card-text class="px-6 pb-6 text-center">
      <v-img
        v-if="!loading && qrImageUrl"
        :src="qrImageUrl"
        max-width="460"
        class="mx-auto"
        cover
      />

      <div v-else-if="loading" class="py-10">
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
const qrImageUrl = ref<string | null>(null)

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

const releaseQr = () => {
  if (qrImageUrl.value) {
    URL.revokeObjectURL(qrImageUrl.value)
    qrImageUrl.value = null
  }
}

const loadQr = async () => {
  loading.value = true

  try {
    releaseQr()
    const blob = await wireguardApi.getPeerQr(peerId.value)
    qrImageUrl.value = URL.createObjectURL(blob)
  } catch (error) {
    notification.showError(api.getErrorMessage(error, 'Не вдалося отримати QR код'))
  } finally {
    loading.value = false
  }
}

onBeforeUnmount(() => {
  releaseQr()
})

await loadQr()
</script>
