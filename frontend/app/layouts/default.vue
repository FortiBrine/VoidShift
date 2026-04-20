<template>
  <v-app>
    <v-navigation-drawer
      v-if="showNavigation"
      :model-value="drawerOpen"
      :temporary="display.smAndDown.value"
      width="260"
      @update:model-value="drawerOpen = $event"
    >
      <v-list nav density="comfortable">
        <v-list-subheader>VoidShift</v-list-subheader>
        <v-list-item
          title="WireGuard"
          prepend-icon="mdi-shield-key-outline"
          to="/wireguard"
          rounded="lg"
        />
      </v-list>

      <template #append>
        <div class="pa-4">
          <v-btn
            block
            variant="outlined"
            prepend-icon="mdi-logout"
            :loading="logoutLoading"
            @click="handleLogout"
          >
            Вийти
          </v-btn>
        </div>
      </template>
    </v-navigation-drawer>

    <v-app-bar v-if="showNavigation" flat>
      <v-app-bar-nav-icon
        v-if="display.smAndDown.value"
        @click="drawerOpen = !drawerOpen"
      />
      <v-app-bar-title>VoidShift</v-app-bar-title>
    </v-app-bar>

    <v-main>
      <v-container class="py-6" max-width="1200">
        <slot />
      </v-container>
    </v-main>

    <v-snackbar
      v-model="notification.visible.value"
      :color="notification.color.value"
      :timeout="notification.timeout"
    >
      {{ notification.message.value }}
    </v-snackbar>
  </v-app>
</template>

<script setup lang="ts">
import { useDisplay } from 'vuetify'

const route = useRoute()
const router = useRouter()
const display = useDisplay()
const auth = useAuth()
const notification = useNotification()

const drawerOpen = ref(true)
const logoutLoading = ref(false)

const showNavigation = computed(() => route.path !== '/login')

watch(
  () => route.path,
  () => {
    if (display.smAndDown.value) {
      drawerOpen.value = false
    }
  },
  { immediate: true },
)

const handleLogout = async () => {
  if (logoutLoading.value) {
    return
  }

  logoutLoading.value = true

  try {
    await auth.logout()
    notification.showSuccess('Сесію завершено')
    await router.push('/login')
  } catch {
    notification.showError('Не вдалося завершити сесію')
  } finally {
    logoutLoading.value = false
  }
}
</script>
