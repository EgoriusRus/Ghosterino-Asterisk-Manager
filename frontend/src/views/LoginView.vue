<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuth } from '@/stores/auth'
import { mdiAccount, mdiLock, mdiPhoneVoip } from '@mdi/js'

const router = useRouter()
const { login, loading, error } = useAuth()

const form = ref({
  username: '',
  password: ''
})

const formValid = ref(true)
const showPassword = ref(false)

const rules = {
  required: (v: string) => !!v || 'Обязательное поле'
}

async function handleLogin() {
  if (!formValid.value) return

  const success = await login({
    username: form.value.username,
    password: form.value.password
  })

  if (success) {
    router.push('/admin/profiles')
  }
}
</script>

<template>
  <v-app>
    <v-main class="d-flex align-center justify-center" style="min-height: 100vh; background: #f8fafc;">
      <v-card width="400" class="pa-4">
        <v-card-text class="text-center pb-0">
          <v-icon :icon="mdiPhoneVoip" size="64" color="primary" class="mb-4" />
          <h1 class="text-h5 font-weight-bold mb-1">Asterisk Manager</h1>
          <p class="text-body-2 text-medium-emphasis">Войдите в систему</p>
        </v-card-text>

        <v-card-text>
          <v-alert
            v-if="error"
            type="error"
            variant="tonal"
            class="mb-4"
            closable
          >
            {{ error }}
          </v-alert>

          <v-form v-model="formValid" @submit.prevent="handleLogin">
            <v-text-field
              v-model="form.username"
              label="Логин"
              :prepend-inner-icon="mdiAccount"
              :rules="[rules.required]"
              variant="outlined"
              class="mb-3"
              autofocus
            />

            <v-text-field
              v-model="form.password"
              label="Пароль"
              :prepend-inner-icon="mdiLock"
              :type="showPassword ? 'text' : 'password'"
              :append-inner-icon="showPassword ? 'mdi-eye-off' : 'mdi-eye'"
              @click:append-inner="showPassword = !showPassword"
              :rules="[rules.required]"
              variant="outlined"
              class="mb-4"
            />

            <v-btn
              type="submit"
              color="primary"
              size="large"
              block
              :loading="loading"
              :disabled="!formValid"
            >
              Войти
            </v-btn>
          </v-form>
        </v-card-text>
      </v-card>
    </v-main>
  </v-app>
</template>
