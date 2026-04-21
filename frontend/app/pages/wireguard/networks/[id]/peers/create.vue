<template>
  <div class="mb-6">
    <v-btn variant="text" prepend-icon="mdi-arrow-left" :to="backPath">Назад до мережі</v-btn>
  </div>

  <v-card max-width="760" rounded="xl" elevation="1">
    <v-card-title class="text-h5 pt-6 px-6">Додати peer</v-card-title>
    <v-card-subtitle class="px-6 pb-2">Створи новий peer для цієї мережі</v-card-subtitle>

    <v-card-text class="px-6 pb-6">
      <v-form @submit.prevent="submitPeer">
        <v-text-field
          v-model="peerIp"
          label="IP peer-а"
          variant="outlined"
          prepend-inner-icon="mdi-ip-network-outline"
          hint="Наприклад: 10.8.0.2"
          persistent-hint
          :rules="peerIpRules"
        />

        <div class="d-flex justify-end mt-4">
          <v-btn color="primary" :loading="loading" @click="submitPeer">Додати peer</v-btn>
        </div>
      </v-form>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
const route = useRoute()
const router = useRouter()
const wireguardApi = useWireguardApi()
const api = useApiClient()
const notification = useNotification()

const loading = ref(false)
const peerIp = ref('')

const networkId = computed(() => Number(route.params.id))
const backPath = computed(() => `/wireguard/networks/${networkId.value}`)

const ipv4Pattern = /^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-5][0-9])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-5][0-9])$/

const peerIpRules = [
  (v: string) => !!v || 'IP обовʼязковий',
  (v: string) => ipv4Pattern.test(v) || 'Некоректний IPv4',
]

if (!Number.isInteger(networkId.value) || networkId.value <= 0) {
  notification.showError('Некоректний ідентифікатор мережі')
  await router.push('/wireguard')
}

const submitPeer = async () => {
  if (loading.value) {
    return
  }

  if (!ipv4Pattern.test(peerIp.value)) {
    notification.showError('Вкажи коректну IPv4 адресу peer-а')
    return
  }

  loading.value = true

  try {
    await wireguardApi.addPeer(networkId.value, {
      allowed_ips: [peerIp.value],
    })

    notification.showSuccess('Peer додано')
    await router.push(backPath.value)
  } catch (error) {
    notification.showError(api.getErrorMessage(error, 'Не вдалося додати peer'))
  } finally {
    loading.value = false
  }
}
</script>
