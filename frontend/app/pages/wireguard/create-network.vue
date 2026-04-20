<template>
  <div class="mb-6">
    <v-btn variant="text" prepend-icon="mdi-arrow-left" to="/wireguard">Назад</v-btn>
  </div>

  <v-card max-width="760" rounded="xl" elevation="1">
    <v-card-title class="text-h5 pt-6 px-6">Створити мережу</v-card-title>
    <v-card-subtitle class="px-6 pb-2">Новий WireGuard-інтерфейс для peer-підключень</v-card-subtitle>

    <v-card-text class="px-6 pb-6">
      <v-form @submit.prevent="submit">
        <v-text-field
          v-model="name"
          label="Назва мережі"
          variant="outlined"
          prepend-inner-icon="mdi-tag-outline"
          :rules="nameRules"
        />

        <v-text-field
          v-model="address"
          label="CIDR мережі"
          variant="outlined"
          prepend-inner-icon="mdi-ip-network-outline"
          hint="Наприклад: 10.8.0.1/24"
          persistent-hint
          :rules="cidrRules"
        />

        <v-text-field
          v-model="listenPort"
          label="Порт підключення"
          variant="outlined"
          prepend-inner-icon="mdi-connection"
          type="number"
          :rules="portRules"
        />

        <div class="d-flex justify-end mt-4">
          <v-btn type="submit" color="primary" :loading="loading">
            Створити
          </v-btn>
        </div>
      </v-form>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
const router = useRouter()
const wireguardApi = useWireguardApi()
const api = useApiClient()
const notification = useNotification()

const loading = ref(false)

const name = ref('')
const address = ref('10.8.0.1/24')
const listenPort = ref('51820')

const cidrPattern = /^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-5][0-9])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-5][0-9])\/(3[0-2]|[12]?[0-9])$/

const nameRules = [
  (v: string) => !!v || 'Назва обовʼязкова',
  (v: string) => v.length >= 4 || 'Мінімум 4 символи',
  (v: string) => v.length <= 100 || 'Максимум 100 символів',
]

const cidrRules = [
  (v: string) => !!v || 'CIDR обовʼязковий',
  (v: string) => cidrPattern.test(v) || 'Некоректний CIDR',
]

const portRules = [
  (v: string) => !!v || 'Порт обовʼязковий',
  (v: string) => {
    const port = Number(v)
    return port >= 1024 && port <= 65535 || 'Порт має бути в межах 1024-65535'
  },
]

const submit = async () => {
  if (loading.value) {
    return
  }

  loading.value = true

  try {
    const created = await wireguardApi.createNetwork({
      name: name.value,
      address: address.value,
      listen_port: Number(listenPort.value),
    })

    notification.showSuccess('Мережу створено')
    await router.push(`/wireguard/networks/${created.id}`)
  } catch (error) {
    notification.showError(api.getErrorMessage(error, 'Не вдалося створити мережу'))
  } finally {
    loading.value = false
  }
}
</script>
