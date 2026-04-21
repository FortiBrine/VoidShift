<template>
  <div class="mb-6 d-flex align-center justify-space-between">
    <v-btn variant="text" prepend-icon="mdi-arrow-left" to="/wireguard">До списку мереж</v-btn>

    <div class="d-flex ga-2">
      <v-btn
        color="primary"
        variant="tonal"
        prepend-icon="mdi-play-circle-outline"
        :loading="upLoading"
        @click="bringUp"
      >
        Up
      </v-btn>

      <v-btn
        color="warning"
        variant="tonal"
        prepend-icon="mdi-pause-circle-outline"
        :loading="downLoading"
        @click="bringDown"
      >
        Down
      </v-btn>
    </div>
  </div>

  <v-card v-if="network" rounded="xl" elevation="1" class="mb-6">
    <v-card-title class="text-h5 pt-6 px-6">Мережа #{{ network.id }}</v-card-title>
    <v-card-text class="px-6 pb-6">
      <v-row>
        <v-col cols="12" md="6">
          <v-list density="compact">
            <v-list-item title="CIDR" :subtitle="network.address" />
            <v-list-item title="Порт" :subtitle="String(network.listen_port)" />
            <v-list-item title="Public key" :subtitle="network.public_key" />
          </v-list>
        </v-col>
      </v-row>
    </v-card-text>
  </v-card>

  <v-card rounded="xl" elevation="1">
    <v-card-title class="d-flex align-center justify-space-between py-4 px-6">
      <span class="text-h6">Peer-и</span>
      <v-btn
        v-if="network"
        color="primary"
        prepend-icon="mdi-plus"
        :to="`/wireguard/networks/${network.id}/peers/create`"
      >
        Додати peer
      </v-btn>
    </v-card-title>

    <v-divider />

    <v-list v-if="network && network.peers.length > 0" lines="three">
      <v-list-item v-for="peer in network.peers" :key="peer.id" :title="`Peer #${peer.id}`" :subtitle="peer.public_key">
        <template #append>
          <div class="d-flex ga-2">
            <v-btn
              v-if="network"
              size="small"
              variant="outlined"
              @click="goToPeerConfig(peer.id)"
            >
              Конфіг
            </v-btn>
            <v-btn size="small" variant="outlined" @click="downloadConfig(peer.id)">Завантажити</v-btn>
            <v-btn
              v-if="network"
              size="small"
              color="primary"
              variant="tonal"
              @click="goToPeerQr(peer.id)"
            >
              QR
            </v-btn>
          </div>
        </template>

        <div class="d-flex ga-2 flex-wrap pt-2">
          <v-chip
            v-for="allowedIp in peer.allowed_ips"
            :key="`${peer.id}-${allowedIp}`"
            size="small"
            variant="tonal"
          >
            {{ allowedIp }}
          </v-chip>
        </div>
      </v-list-item>
    </v-list>

    <v-card-text v-else-if="!loading" class="py-10 text-center text-medium-emphasis">
      Peer-ів поки немає.
    </v-card-text>

    <v-card-text v-if="loading" class="py-10 text-center">
      <v-progress-circular indeterminate color="primary" />
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import type { NetworkDetails } from '~/composables/useWireguardApi'

const route = useRoute()
const router = useRouter()
const wireguardApi = useWireguardApi()
const api = useApiClient()
const notification = useNotification()

const loading = ref(true)
const upLoading = ref(false)
const downLoading = ref(false)

const network = ref<NetworkDetails | null>(null)

const networkId = computed(() => Number(route.params.id))

const loadNetwork = async () => {
  if (!Number.isInteger(networkId.value) || networkId.value <= 0) {
    notification.showError('Некоректний ідентифікатор мережі')
    await router.push('/wireguard')
    return
  }

  loading.value = true

  try {
    network.value = await wireguardApi.getNetwork(networkId.value)
  } catch (error) {
    notification.showError(api.getErrorMessage(error, 'Не вдалося завантажити мережу'))
  } finally {
    loading.value = false
  }
}

const bringUp = async () => {
  if (upLoading.value) {
    return
  }

  upLoading.value = true

  try {
    await wireguardApi.networkUp(networkId.value)
    notification.showSuccess('Мережу піднято')
  } catch (error) {
    notification.showError(api.getErrorMessage(error, 'Не вдалося підняти мережу'))
  } finally {
    upLoading.value = false
  }
}

const bringDown = async () => {
  if (downLoading.value) {
    return
  }

  downLoading.value = true

  try {
    await wireguardApi.networkDown(networkId.value)
    notification.showSuccess('Мережу зупинено')
  } catch (error) {
    notification.showError(api.getErrorMessage(error, 'Не вдалося зупинити мережу'))
  } finally {
    downLoading.value = false
  }
}

const downloadConfig = (peerId: number) => {
  window.open(`/api/vpn/wireguard/peers/${peerId}/config/download`, '_blank', 'noopener')
}

const goToPeerConfig = async (peerId: number) => {
  if (!network.value) {
    return
  }

  await router.push({
    path: `/wireguard/peers/${peerId}/config`,
    query: { networkId: String(network.value.id) },
  })
}

const goToPeerQr = async (peerId: number) => {
  if (!network.value) {
    return
  }

  await router.push({
    path: `/wireguard/peers/${peerId}/qr`,
    query: { networkId: String(network.value.id) },
  })
}

await loadNetwork()
</script>
