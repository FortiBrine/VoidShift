<template>
  <v-container class="fill-height d-flex align-center justify-center" max-width="420">
    <v-card width="100%" rounded="xl" elevation="2">
      <v-card-title class="text-h5 pt-6 pb-2 px-6">Вхід у VoidShift</v-card-title>
      <v-card-subtitle class="px-6 pb-2">
        Self-hosted панель керування WireGuard
      </v-card-subtitle>

      <v-card-text class="px-6 pb-2">
        <v-form @submit.prevent="submit">
          <v-text-field
            v-model="username"
            label="Логін"
            variant="outlined"
            prepend-inner-icon="mdi-account-outline"
            autocomplete="username"
            :rules="usernameRules"
          />

          <v-text-field
            v-model="password"
            label="Пароль"
            variant="outlined"
            prepend-inner-icon="mdi-lock-outline"
            :append-inner-icon="showPassword ? 'mdi-eye-off-outline' : 'mdi-eye-outline'"
            :type="showPassword ? 'text' : 'password'"
            autocomplete="current-password"
            :rules="passwordRules"
            @click:append-inner="showPassword = !showPassword"
          />

          <v-btn
            type="submit"
            block
            size="large"
            color="primary"
            :loading="loading"
          >
            Увійти
          </v-btn>
        </v-form>
      </v-card-text>
    </v-card>
  </v-container>
</template>

<script setup lang="ts">
const router = useRouter()
const auth = useAuth()
const api = useApiClient()
const notification = useNotification()

const loading = ref(false)
const showPassword = ref(false)

const username = ref('')
const password = ref('')

const usernameRules = [
  (v: string) => !!v || 'Логін обовʼязковий',
  (v: string) => v.length >= 4 || 'Мінімум 4 символи',
  (v: string) => v.length <= 30 || 'Максимум 30 символів',
  (v: string) => /^[a-zA-Z0-9]+$/.test(v) || 'Тільки букви й цифри',
]

const passwordRules = [
  (v: string) => !!v || 'Пароль обовʼязковий',
  (v: string) => v.length >= 8 || 'Мінімум 8 символів',
  (v: string) => v.length <= 40 || 'Максимум 40 символів',
]

const submit = async () => {
  if (loading.value) {
    return
  }

  loading.value = true

  try {
    await auth.login(username.value, password.value)
    notification.showSuccess('Вхід виконано')
    await router.push('/wireguard')
  } catch (error) {
    notification.showError(api.getErrorMessage(error, 'Не вдалося увійти'))
  } finally {
    loading.value = false
  }
}
</script>
